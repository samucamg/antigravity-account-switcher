package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	_ "time/tzdata"

	"github.com/samucamg/antigravity-account-switcher/internal/config"
	"github.com/samucamg/antigravity-account-switcher/internal/domain"
	"github.com/samucamg/antigravity-account-switcher/internal/oauth"
	"github.com/samucamg/antigravity-account-switcher/internal/quota"
	"github.com/samucamg/antigravity-account-switcher/internal/tunnel"
)

// FallbackConfigSetter defines an interface for dynamically updating model fallback settings at runtime.
type FallbackConfigSetter interface {
	SetFallbackConfig(primary, secondary string, enabled bool)
}

// APIHandler implements the REST endpoints and SSE real-time streaming for the switcher.
type APIHandler struct {
	accountRepo          domain.AccountRepository
	quotaRepo            domain.QuotaRepository
	metricsService       domain.MetricsService
	broadcaster          domain.EventBroadcaster
	eventRepo            domain.EventRepository
	oauthEngine          oauth.OAuthEngine
	poller               QuotaPoller
	startTime            time.Time
	version              string
	cfgMu                sync.RWMutex
	appConfig            *config.Config
	fallbackConfigSetter FallbackConfigSetter
	tunnelManager        *tunnel.Manager
}

// SetConfig sets the configuration pointer for APIHandler.
func (a *APIHandler) SetConfig(c *config.Config) {
	a.cfgMu.Lock()
	defer a.cfgMu.Unlock()
	a.appConfig = c
}

// SetFallbackConfigSetter sets the dynamic fallback setter for live proxy updates.
func (a *APIHandler) SetFallbackConfigSetter(s FallbackConfigSetter) {
	a.cfgMu.Lock()
	defer a.cfgMu.Unlock()
	a.fallbackConfigSetter = s
}

// SetTunnelManager sets the cloudflare tunnel manager.
func (a *APIHandler) SetTunnelManager(t *tunnel.Manager) {
	a.cfgMu.Lock()
	defer a.cfgMu.Unlock()
	a.tunnelManager = t
}

// QuotaPoller defines interface for triggering quota polling passes.
type QuotaPoller interface {
	PollOnce(ctx context.Context) error
}

// SetPoller assigns the quota poller.
func (a *APIHandler) SetPoller(p QuotaPoller) {
	a.poller = p
}

// HandleQuotaRefresh serves POST /api/quota/refresh.
func (a *APIHandler) HandleQuotaRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if a.poller != nil {
		_ = a.poller.PollOnce(r.Context())
	}

	a.listAccounts(w, r)
}

// NewAPIHandler constructs an APIHandler with provided dependencies.
func NewAPIHandler(
	accountRepo domain.AccountRepository,
	quotaRepo domain.QuotaRepository,
	metricsService domain.MetricsService,
	broadcaster domain.EventBroadcaster,
	eventRepo domain.EventRepository,
	oauthEngine oauth.OAuthEngine,
	version string,
) *APIHandler {
	if version == "" {
		version = "0.1.0-dev"
	}
	return &APIHandler{
		accountRepo:    accountRepo,
		quotaRepo:      quotaRepo,
		metricsService: metricsService,
		broadcaster:    broadcaster,
		eventRepo:      eventRepo,
		oauthEngine:    oauthEngine,
		startTime:      time.Now(),
		version:        version,
	}
}

// StatusResponse represents server health and active account info.
type StatusResponse struct {
	Status        string          `json:"status"`
	Version       string          `json:"version"`
	UptimeSeconds int64           `json:"uptime_seconds"`
	ActiveAccount *domain.Account `json:"active_account,omitempty"`
	TotalAccounts int             `json:"total_accounts"`
	Timestamp     time.Time       `json:"timestamp"`
}

// HandleStatus serves GET /api/status.
func (a *APIHandler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	var active *domain.Account
	total := 0

	if a.accountRepo != nil {
		active, _ = a.accountRepo.GetActive(ctx)
		if list, err := a.accountRepo.List(ctx); err == nil {
			total = len(list)
		}
	}

	resp := StatusResponse{
		Status:        "ok",
		Version:       a.version,
		UptimeSeconds: int64(time.Since(a.startTime).Seconds()),
		ActiveAccount: active,
		TotalAccounts: total,
		Timestamp:     time.Now().UTC(),
	}

	writeJSON(w, http.StatusOK, resp)
}

// AccountWithBuckets encapsulates an Account and its associated quota buckets.
type AccountWithBuckets struct {
	*domain.Account
	Buckets []*domain.QuotaBucket `json:"buckets"`
}

// HandleAccounts serves GET /api/accounts and DELETE /api/accounts/{id}.
func (a *APIHandler) HandleAccounts(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/accounts")
	path = strings.Trim(path, "/")

	if path == "" {
		if r.Method == http.MethodGet {
			a.listAccounts(w, r)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(path, "/")
	accountID := parts[0]

	if len(parts) == 2 && parts[1] == "select" {
		if r.Method == http.MethodPost {
			a.selectAccount(w, r, accountID)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if len(parts) == 1 {
		switch r.Method {
		case http.MethodGet:
			a.getAccount(w, r, accountID)
			return
		case http.MethodDelete:
			a.deleteAccount(w, r, accountID)
			return
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}

	http.NotFound(w, r)
}

func (a *APIHandler) listAccounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if a.accountRepo == nil {
		writeJSON(w, http.StatusOK, []*AccountWithBuckets{})
		return
	}

	accounts, err := a.accountRepo.List(ctx)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, "failed to list accounts", err)
		return
	}

	var allBuckets map[string][]*domain.QuotaBucket
	if a.quotaRepo != nil {
		allBuckets, _ = a.quotaRepo.ListAll(ctx)
	}

	result := make([]*AccountWithBuckets, 0, len(accounts))
	for _, acc := range accounts {
		var b []*domain.QuotaBucket
		if allBuckets != nil {
			b = allBuckets[acc.ID]
		}
		if b == nil {
			b = []*domain.QuotaBucket{}
		}
		result = append(result, &AccountWithBuckets{
			Account: acc,
			Buckets: b,
		})
	}

	writeJSON(w, http.StatusOK, result)
}

func (a *APIHandler) getAccount(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	if a.accountRepo == nil {
		http.NotFound(w, r)
		return
	}

	acc, err := a.accountRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			http.NotFound(w, r)
			return
		}
		writeErrorJSON(w, http.StatusInternalServerError, "failed to get account", err)
		return
	}

	var buckets []*domain.QuotaBucket
	if a.quotaRepo != nil {
		buckets, _ = a.quotaRepo.GetByAccountID(ctx, id)
	}
	if buckets == nil {
		buckets = []*domain.QuotaBucket{}
	}

	writeJSON(w, http.StatusOK, &AccountWithBuckets{
		Account: acc,
		Buckets: buckets,
	})
}

func (a *APIHandler) selectAccount(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	if a.accountRepo == nil {
		writeErrorJSON(w, http.StatusServiceUnavailable, "account repository unavailable", nil)
		return
	}

	acc, err := a.accountRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			writeErrorJSON(w, http.StatusNotFound, "account not found", err)
			return
		}
		writeErrorJSON(w, http.StatusInternalServerError, "failed to lookup account", err)
		return
	}

	if err := a.accountRepo.SetActive(ctx, id); err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, "failed to activate account", err)
		return
	}

	evt := &domain.ProxyEvent{
		Type:      domain.EventTypeAccountSwitched,
		AccountID: id,
		Message:   fmt.Sprintf("Account %s (%s) set as active via dashboard", acc.Email, id),
		Timestamp: time.Now().UTC(),
	}
	if a.broadcaster != nil {
		a.broadcaster.Broadcast(evt)
	}
	if a.eventRepo != nil {
		_ = a.eventRepo.Record(ctx, evt)
	}

	if a.poller != nil {
		go func() {
			_ = a.poller.PollOnce(context.Background())
		}()
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":    true,
		"account_id": id,
		"email":      acc.Email,
	})
}

func (a *APIHandler) deleteAccount(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	if a.accountRepo == nil {
		writeErrorJSON(w, http.StatusServiceUnavailable, "account repository unavailable", nil)
		return
	}

	acc, err := a.accountRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			writeErrorJSON(w, http.StatusNotFound, "account not found", err)
			return
		}
		writeErrorJSON(w, http.StatusInternalServerError, "failed to lookup account", err)
		return
	}

	wasActive := acc.IsActive

	if err := a.accountRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			writeErrorJSON(w, http.StatusNotFound, "account not found", err)
			return
		}
		writeErrorJSON(w, http.StatusInternalServerError, "failed to delete account", err)
		return
	}

	if a.quotaRepo != nil {
		_ = a.quotaRepo.DeleteByAccountID(ctx, id)
	}

	// If the deleted account was active, auto-promote next available account
	if wasActive {
		if next, nextErr := a.accountRepo.GetNextAvailable(ctx, ""); nextErr == nil && next != nil {
			_ = a.accountRepo.SetActive(ctx, next.ID)
			if a.broadcaster != nil {
				a.broadcaster.Broadcast(&domain.ProxyEvent{
					Type:      domain.EventTypeAccountSwitched,
					AccountID: next.ID,
					Message:   fmt.Sprintf("Account %s automatically promoted to active following deletion of %s", next.Email, acc.Email),
					Timestamp: time.Now().UTC(),
				})
			}
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":    true,
		"account_id": id,
	})
}

// HandleMetrics serves GET /api/metrics.
func (a *APIHandler) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	if a.metricsService == nil {
		writeJSON(w, http.StatusOK, &domain.MetricsDashboardPayload{
			Summary: domain.GlobalDashboardSummary{
				Today:     &domain.AggregatedMetrics{},
				ThisWeek:  &domain.AggregatedMetrics{},
				ThisMonth: &domain.AggregatedMetrics{},
				AllTime:   &domain.AggregatedMetrics{},
			},
			ByAccount: []*domain.AccountMetricsSummary{},
			Timeline:  []*domain.DailyTokenUsage{},
		})
		return
	}

	accountID := r.URL.Query().Get("account_id")
	periodParam := r.URL.Query().Get("period")
	tzParam := r.URL.Query().Get("tz")
	tzOffsetParam := r.URL.Query().Get("tz_offset")
	loc := parseLocation(tzParam, tzOffsetParam)

	if accountID != "" {
		norm := domain.PeriodLifetime
		if periodParam != "" {
			norm = domain.MetricPeriod(strings.ToLower(periodParam))
		}
		summary, err := a.metricsService.GetSummary(ctx, accountID, norm)
		if err != nil {
			writeErrorJSON(w, http.StatusInternalServerError, "failed to get account summary", err)
			return
		}
		writeJSON(w, http.StatusOK, summary)
		return
	}

	payload, err := a.metricsService.GetDashboardPayloadWithLocation(ctx, 14, loc)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, "failed to compute metrics dashboard payload", err)
		return
	}

	writeJSON(w, http.StatusOK, payload)
}

// parseLocation resolves a *time.Location from either an IANA timezone name or a numeric minute offset.
func parseLocation(tzParam, tzOffsetParam string) *time.Location {
	if tzParam != "" {
		if loc, err := time.LoadLocation(tzParam); err == nil {
			return loc
		}
	}
	if tzOffsetParam != "" {
		// tz_offset in minutes: -new Date().getTimezoneOffset()
		// e.g. -180 for UTC-3, +540 for UTC+9
		if offsetMinutes, err := strconv.Atoi(tzOffsetParam); err == nil {
			hours := offsetMinutes / 60
			mins := int(math.Abs(float64(offsetMinutes % 60)))
			name := fmt.Sprintf("UTC%+03d:%02d", hours, mins)
			return time.FixedZone(name, offsetMinutes*60)
		}
	}
	return time.UTC
}

// HandleEvents serves GET /api/events as a real-time SSE stream.
func (a *APIHandler) HandleEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported by client connection", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-transform")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	ctx := r.Context()

	// 1. Send recent historical events if event repository is present
	if a.eventRepo != nil {
		recent, err := a.eventRepo.ListRecent(ctx, 30)
		if err == nil && len(recent) > 0 {
			// ListRecent returns newest first; reverse for chronological playback
			for i := len(recent) - 1; i >= 0; i-- {
				evt := recent[i]
				if data, err := json.Marshal(evt); err == nil {
					_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
				}
			}
			flusher.Flush()
		}
	}

	// 2. Stream real-time events via broadcaster
	if a.broadcaster == nil {
		<-ctx.Done()
		return
	}

	eventChan, unsubscribe := a.broadcaster.Subscribe()
	defer unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return
		case evt, ok := <-eventChan:
			if !ok {
				return
			}
			data, err := json.Marshal(evt)
			if err != nil {
				continue
			}
			_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}

// HandleOAuthStart serves POST/GET /oauth/start.
func (a *APIHandler) HandleOAuthStart(w http.ResponseWriter, r *http.Request) {
	if a.oauthEngine == nil {
		writeErrorJSON(w, http.StatusNotImplemented, "OAuth2 engine not configured", nil)
		return
	}

	urlChan := make(chan string, 1)

	// Non-blocking trigger of loopback flow
	go func() {
		_, err := a.oauthEngine.StartLoopbackFlow(context.Background(), nil, func(authURL string) {
			select {
			case urlChan <- authURL:
			default:
			}
			if a.broadcaster != nil {
				a.broadcaster.Broadcast(&domain.ProxyEvent{
					Type:      domain.EventType("oauth_started"),
					Message:   fmt.Sprintf("OAuth authorization flow initiated: %s", authURL),
					Timestamp: time.Now().UTC(),
				})
			}
		})
		if err != nil && a.broadcaster != nil {
			a.broadcaster.Broadcast(&domain.ProxyEvent{
				Type:      domain.EventTypeError,
				Message:   fmt.Sprintf("OAuth flow failed: %v", err),
				Timestamp: time.Now().UTC(),
			})
		}
	}()

	var generatedAuthURL string
	select {
	case generatedAuthURL = <-urlChan:
	case <-time.After(1 * time.Second):
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":   "started",
		"message":  "OAuth loopback flow initiated in browser",
		"auth_url": generatedAuthURL,
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func writeErrorJSON(w http.ResponseWriter, statusCode int, message string, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	detail := ""
	if err != nil {
		detail = err.Error()
	}
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]any{
			"code":    statusCode,
			"message": message,
			"detail":  detail,
		},
	})
}

// ConfigResponse represents the response payload for GET /api/config.
type ConfigResponse struct {
	ModelPrimary             string `json:"model_primary"`
	ModelSecondary           string `json:"model_secondary"`
	FallbackSecondaryEnabled bool   `json:"fallback_secondary_enabled"`
	CloudflareTunnelToken    string `json:"cloudflare_tunnel_token,omitempty"`
	RemoteAuthEnabled        bool   `json:"remote_auth_enabled"`
}

// ConfigUpdateRequest represents the payload for POST /api/config.
type ConfigUpdateRequest struct {
	ModelPrimary             *string `json:"model_primary,omitempty"`
	ModelSecondary           *string `json:"model_secondary,omitempty"`
	FallbackSecondaryEnabled *bool   `json:"fallback_secondary_enabled,omitempty"`
	CloudflareTunnelToken    *string `json:"cloudflare_tunnel_token,omitempty"`
	RemoteAuthToken          *string `json:"remote_auth_token,omitempty"`
}

// HandleConfig serves GET and POST /api/config.
func (a *APIHandler) HandleConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.getConfig(w, r)
	case http.MethodPost, http.MethodPut:
		a.updateConfig(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *APIHandler) getConfig(w http.ResponseWriter, _ *http.Request) {
	a.cfgMu.RLock()
	var primary, secondary, tunnelToken string
	var enabled, remoteAuth bool
	if a.appConfig != nil {
		primary = a.appConfig.ModelPrimary
		secondary = a.appConfig.ModelSecondary
		enabled = a.appConfig.FallbackSecondaryEnabled
		tunnelToken = a.appConfig.CloudflareTunnelToken
		remoteAuth = a.appConfig.RemoteAuthToken != ""
	}
	a.cfgMu.RUnlock()

	writeJSON(w, http.StatusOK, ConfigResponse{
		ModelPrimary:             primary,
		ModelSecondary:           secondary,
		FallbackSecondaryEnabled: enabled,
		CloudflareTunnelToken:    tunnelToken,
		RemoteAuthEnabled:        remoteAuth,
	})
}

func (a *APIHandler) updateConfig(w http.ResponseWriter, r *http.Request) {
	var req ConfigUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "invalid JSON payload", err)
		return
	}

	a.cfgMu.Lock()
	if a.appConfig == nil {
		a.appConfig = config.DefaultConfig()
	}
	currentCfg := a.appConfig

	if req.ModelPrimary != nil {
		if trimmed := strings.TrimSpace(*req.ModelPrimary); trimmed != "" {
			currentCfg.ModelPrimary = trimmed
		}
	}
	if req.ModelSecondary != nil {
		if trimmed := strings.TrimSpace(*req.ModelSecondary); trimmed != "" {
			currentCfg.ModelSecondary = trimmed
		}
	}
	if req.FallbackSecondaryEnabled != nil {
		currentCfg.FallbackSecondaryEnabled = *req.FallbackSecondaryEnabled
	}
	if req.CloudflareTunnelToken != nil {
		currentCfg.CloudflareTunnelToken = strings.TrimSpace(*req.CloudflareTunnelToken)
	}
	if req.RemoteAuthToken != nil {
		currentCfg.RemoteAuthToken = strings.TrimSpace(*req.RemoteAuthToken)
	}

	_ = config.Save(currentCfg)

	if a.fallbackConfigSetter != nil {
		a.fallbackConfigSetter.SetFallbackConfig(
			currentCfg.ModelPrimary,
			currentCfg.ModelSecondary,
			currentCfg.FallbackSecondaryEnabled,
		)
	}

	resp := ConfigResponse{
		ModelPrimary:             currentCfg.ModelPrimary,
		ModelSecondary:           currentCfg.ModelSecondary,
		FallbackSecondaryEnabled: currentCfg.FallbackSecondaryEnabled,
		CloudflareTunnelToken:    currentCfg.CloudflareTunnelToken,
		RemoteAuthEnabled:        currentCfg.RemoteAuthToken != "",
	}
	broadcaster := a.broadcaster
	a.cfgMu.Unlock()

	if broadcaster != nil {
		broadcaster.Broadcast(&domain.ProxyEvent{
			Type:    domain.EventTypeAccountSwitched,
			Message: fmt.Sprintf("Model fallback updated: Primary=%s, Secondary=%s, Enabled=%t", resp.ModelPrimary, resp.ModelSecondary, resp.FallbackSecondaryEnabled),
			Details: map[string]any{
				"model_primary":              resp.ModelPrimary,
				"model_secondary":            resp.ModelSecondary,
				"fallback_secondary_enabled": resp.FallbackSecondaryEnabled,
			},
			Timestamp: time.Now().UTC(),
		})
	}

	writeJSON(w, http.StatusOK, resp)
}

// ModelsResponse models the JSON payload for GET /api/models.
type ModelsResponse struct {
	Models []*domain.ModelInfo `json:"models"`
	Source string              `json:"source"`
}

// HandleModels serves GET /api/models.
func (a *APIHandler) HandleModels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	models, source := a.discoverModels(ctx)

	writeJSON(w, http.StatusOK, ModelsResponse{
		Models: models,
		Source: source,
	})
}

func (a *APIHandler) discoverModels(ctx context.Context) ([]*domain.ModelInfo, string) {
	// 1. Try querying running language_server on localhost
	if lsModels, err := quota.QueryAvailableModels(ctx); err == nil && len(lsModels) > 0 {
		return a.ensureConfiguredModelsPresent(lsModels), "language_server"
	}

	// 2. Try querying Cloud Code PA directly if active account token exists
	if a.accountRepo != nil {
		if activeAcc, err := a.accountRepo.GetActive(ctx); err == nil && activeAcc != nil && activeAcc.AccessToken != "" {
			if ccModels, err := quota.FetchAvailableModelsFromCloudCode(ctx, activeAcc.AccessToken); err == nil && len(ccModels) > 0 {
				return a.ensureConfiguredModelsPresent(ccModels), "cloud_code_pa"
			}
		}
	}

	// 3. Fall back to standard curated built-in model matrix
	return a.ensureConfiguredModelsPresent(standardFallbackModels()), "standard_catalog"
}

func (a *APIHandler) ensureConfiguredModelsPresent(discovered []*domain.ModelInfo) []*domain.ModelInfo {
	a.cfgMu.RLock()
	var primary, secondary string
	if a.appConfig != nil {
		primary = a.appConfig.ModelPrimary
		secondary = a.appConfig.ModelSecondary
	}
	a.cfgMu.RUnlock()

	res := make([]*domain.ModelInfo, len(discovered))
	copy(res, discovered)

	hasPrimary := false
	hasSecondary := false
	for _, m := range res {
		if m.ID == primary {
			hasPrimary = true
		}
		if m.ID == secondary {
			hasSecondary = true
		}
	}

	if primary != "" && !hasPrimary {
		res = append(res, &domain.ModelInfo{
			ID:          primary,
			DisplayName: primary,
			Category:    "gemini",
			Recommended: true,
		})
	}
	if secondary != "" && !hasSecondary {
		res = append(res, &domain.ModelInfo{
			ID:          secondary,
			DisplayName: secondary,
			Category:    "claude_gpt",
			Recommended: false,
		})
	}
	return res
}

func standardFallbackModels() []*domain.ModelInfo {
	return []*domain.ModelInfo{
		{ID: "gemini-2.5-pro", DisplayName: "Gemini 2.5 Pro", Category: "gemini", Recommended: true},
		{ID: "gemini-2.5-flash", DisplayName: "Gemini 2.5 Flash", Category: "gemini", Recommended: true},
		{ID: "gemini-2.0-pro-exp", DisplayName: "Gemini 2.0 Pro Exp", Category: "gemini", Recommended: false},
		{ID: "gemini-2.0-flash", DisplayName: "Gemini 2.0 Flash", Category: "gemini", Recommended: false},
		{ID: "claude-3-7-sonnet", DisplayName: "Claude 3.7 Sonnet", Category: "claude_gpt", Recommended: true},
		{ID: "claude-3-5-sonnet", DisplayName: "Claude 3.5 Sonnet", Category: "claude_gpt", Recommended: true},
		{ID: "claude-3-5-haiku", DisplayName: "Claude 3.5 Haiku", Category: "claude_gpt", Recommended: false},
	}
}

// HandleTunnelStatus handles GET /api/tunnel/status
func (a *APIHandler) HandleTunnelStatus(w http.ResponseWriter, r *http.Request) {
	if a.tunnelManager == nil {
		writeJSON(w, http.StatusOK, tunnel.Status{Active: false, Mode: tunnel.ModeNone})
		return
	}
	writeJSON(w, http.StatusOK, a.tunnelManager.GetStatus())
}

// HandleTunnelStart handles POST /api/tunnel/start
func (a *APIHandler) HandleTunnelStart(w http.ResponseWriter, r *http.Request) {
	if a.tunnelManager == nil {
		writeErrorJSON(w, http.StatusInternalServerError, "tunnel manager not initialized", nil)
		return
	}
	var req struct {
		Type  string `json:"type"`
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	port := 8080
	a.cfgMu.RLock()
	if a.appConfig != nil && a.appConfig.Port > 0 {
		port = a.appConfig.Port
	}
	a.cfgMu.RUnlock()

	var err error
	if req.Type == "zero_trust" {
		token := req.Token
		if token == "" {
			a.cfgMu.RLock()
			if a.appConfig != nil {
				token = a.appConfig.CloudflareTunnelToken
			}
			a.cfgMu.RUnlock()
		}
		if token == "" {
			writeErrorJSON(w, http.StatusBadRequest, "zero trust tunnel token is required", nil)
			return
		}
		err = a.tunnelManager.StartTokenTunnel(token)
	} else {
		err = a.tunnelManager.StartQuickTunnel(port)
	}

	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, "failed to start tunnel", err)
		return
	}

	writeJSON(w, http.StatusOK, a.tunnelManager.GetStatus())
}

// HandleTunnelStop handles POST /api/tunnel/stop
func (a *APIHandler) HandleTunnelStop(w http.ResponseWriter, r *http.Request) {
	if a.tunnelManager == nil {
		writeJSON(w, http.StatusOK, map[string]string{"status": "stopped"})
		return
	}
	if err := a.tunnelManager.Stop(); err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, "failed to stop tunnel", err)
		return
	}
	writeJSON(w, http.StatusOK, a.tunnelManager.GetStatus())
}

// CheckAuthorized checks if incoming remote request is authenticated.
func (a *APIHandler) CheckAuthorized(r *http.Request) bool {
	a.cfgMu.RLock()
	expectedToken := ""
	if a.appConfig != nil {
		expectedToken = a.appConfig.RemoteAuthToken
	}
	a.cfgMu.RUnlock()

	if expectedToken == "" {
		return true // No password/token set, free access
	}

	// 1. Check Bearer token
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == expectedToken {
			return true
		}
	}

	// 2. Check query parameter ?token=...
	if r.URL.Query().Get("token") == expectedToken {
		return true
	}

	return false
}
