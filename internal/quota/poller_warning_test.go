package quota

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/samucamg/antigravity-account-switcher/internal/domain"
	"github.com/samucamg/antigravity-account-switcher/internal/store/sqlite"
)

// quotaServer returns an httptest server whose quota summary endpoint responds
// with a single bucket whose remaining fraction is `remaining`.
func quotaServer(t *testing.T, remaining float64) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/v1internal:retrieveUserQuotaSummary", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body := fmt.Sprintf(`{"buckets":[{"bucketId":"claude-gpt-5h","displayName":"Claude and GPT models (5h)","window":"5h","remainingFraction":%.2f,"remainingAmount":100,"resetTime":"2030-01-01T00:00:00Z"}]}`, remaining)
		_, _ = w.Write([]byte(body))
	})
	return httptest.NewServer(mux)
}

func setupAccountForPoll(t *testing.T, accRepo *sqlite.AccountRepository, id, email string) *domain.Account {
	t.Helper()
	// NOTE: is_active=false because the SQLite schema enforces a single active
	// account (UNIQUE on is_active). The poller fetches remote quota buckets
	// regardless when a custom BaseURL is used.
	acc := &domain.Account{
		ID:          id,
		Email:       email,
		AccessToken: "token-" + id,
		TokenExpiry: time.Now().Add(2 * time.Hour),
		Status:      domain.AccountStatusActive,
		IsActive:    false,
	}
	if err := accRepo.Create(context.Background(), acc); err != nil {
		t.Fatalf("failed to create account: %v", err)
	}
	return acc
}

func TestPoller_EmitsQuotaWarning_AboveThreshold(t *testing.T) {
	server := quotaServer(t, 0.10) // 90% usage → above 80% default warning threshold
	defer server.Close()

	_, accRepo, quotaRepo := setupTestStore(t)
	acc := setupAccountForPoll(t, accRepo, "acc-warn", "warn@example.com")

	bc := &mockBroadcaster{}
	p, err := NewPoller(accRepo, quotaRepo,
		WithBaseURL(server.URL),
		WithEventBroadcaster(bc),
	)
	if err != nil {
		t.Fatalf("NewPoller: %v", err)
	}

	if err := p.pollAccount(context.Background(), acc, time.Now().UTC()); err != nil {
		t.Fatalf("pollAccount failed: %v", err)
	}

	events := bc.getEvents()
	if len(events) == 0 {
		t.Fatal("expected at least one quota warning event, got none")
	}

	var foundWarning bool
	for _, ev := range events {
		if ev.Type == domain.EventTypeQuotaWarning {
			foundWarning = true
			if ev.AccountID != "acc-warn" {
				t.Errorf("expected AccountID acc-warn, got %s", ev.AccountID)
			}
			if ev.Details["usage_pct"] == nil {
				t.Error("expected usage_pct in event details")
			}
		}
	}
	if !foundWarning {
		t.Error("expected EventTypeQuotaWarning to be broadcast when usage is above threshold")
	}
}

func TestPoller_NoQuotaWarning_BelowThreshold(t *testing.T) {
	server := quotaServer(t, 0.95) // 5% usage → below 80% default warning threshold
	defer server.Close()

	_, accRepo, quotaRepo := setupTestStore(t)
	acc := setupAccountForPoll(t, accRepo, "acc-quiet", "quiet@example.com")

	bc := &mockBroadcaster{}
	p, err := NewPoller(accRepo, quotaRepo,
		WithBaseURL(server.URL),
		WithEventBroadcaster(bc),
	)
	if err != nil {
		t.Fatalf("NewPoller: %v", err)
	}

	if err := p.pollAccount(context.Background(), acc, time.Now().UTC()); err != nil {
		t.Fatalf("pollAccount failed: %v", err)
	}

	for _, ev := range bc.getEvents() {
		if ev.Type == domain.EventTypeQuotaWarning {
			t.Errorf("unexpected quota warning event emitted at low usage: %+v", ev)
		}
	}
}

func TestPoller_QuotaWarning_RespectsCustomThreshold(t *testing.T) {
	server := quotaServer(t, 0.20) // 80% usage
	defer server.Close()

	_, accRepo, quotaRepo := setupTestStore(t)
	acc := setupAccountForPoll(t, accRepo, "acc-custom", "custom@example.com")

	bc := &mockBroadcaster{}
	p, err := NewPoller(accRepo, quotaRepo,
		WithBaseURL(server.URL),
		WithEventBroadcaster(bc),
		WithQuotaWarningThreshold(0.85), // stricter: 80% usage must NOT warn
	)
	if err != nil {
		t.Fatalf("NewPoller: %v", err)
	}

	if err := p.pollAccount(context.Background(), acc, time.Now().UTC()); err != nil {
		t.Fatalf("pollAccount failed: %v", err)
	}

	for _, ev := range bc.getEvents() {
		if ev.Type == domain.EventTypeQuotaWarning {
			t.Errorf("expected no warning with 80%% usage under 0.85 threshold, got %+v", ev)
		}
	}

	// Same account at 90% usage (remaining 0.10) must warn under 0.85 threshold
	server2 := quotaServer(t, 0.10)
	defer server2.Close()
	acc2 := setupAccountForPoll(t, accRepo, "acc-custom2", "custom2@example.com")

	bc2 := &mockBroadcaster{}
	p2, err := NewPoller(accRepo, quotaRepo,
		WithBaseURL(server2.URL),
		WithEventBroadcaster(bc2),
		WithQuotaWarningThreshold(0.85),
	)
	if err != nil {
		t.Fatalf("NewPoller: %v", err)
	}

	if err := p2.pollAccount(context.Background(), acc2, time.Now().UTC()); err != nil {
		t.Fatalf("pollAccount failed: %v", err)
	}

	var found bool
	for _, ev := range bc2.getEvents() {
		if ev.Type == domain.EventTypeQuotaWarning {
			found = true
		}
	}
	if !found {
		t.Error("expected quota warning at 90% usage under 0.85 threshold")
	}
}