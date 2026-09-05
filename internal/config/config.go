package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	DefaultPort        = 8080
	DefaultUpstreamURL = "https://daily-cloudcode-pa.googleapis.com"
	DefaultInterval    = "5m"
	DefaultQuotaWarningThreshold = 0.80 // 80% usage threshold for alert
	DefaultQuotaSwitchThreshold  = 0.85 // 85% usage threshold for proactive account rotation
	DefaultModelPrimary          = "gemini-2.5-pro"
	DefaultModelSecondary        = "gemini-2.5-flash"
	DefaultFallbackSecondaryEnabled = false
	ConfigFileName     = "config.json"
	DefaultDBFileName  = "accounts.db"
)

// Config holds persistent user configuration.
type Config struct {
	Port                  int     `json:"port"`
	DBPath                string  `json:"db_path"`
	AntigravityBin        string  `json:"antigravity_bin"`
	UpstreamURL           string  `json:"upstream_url"`
	QuotaInterval         string  `json:"quota_interval"`
	QuotaWarningThreshold float64 `json:"quota_warning_threshold"`
	QuotaSwitchThreshold  float64 `json:"quota_switch_threshold"`
	OpenBrowser           bool    `json:"open_browser"`
	ModelPrimary          string  `json:"model_primary,omitempty"`
	ModelSecondary        string  `json:"model_secondary,omitempty"`
	FallbackSecondaryEnabled bool `json:"fallback_secondary_enabled"`
	CloudflareTunnelToken string  `json:"cloudflare_tunnel_token,omitempty"`
	RemoteAuthToken       string  `json:"remote_auth_token,omitempty"`
}

// ConfigDir returns the default configuration directory (~/.config/antigravity-account-switcher).
func ConfigDir() string {
	if custom := os.Getenv("ANTIGRAVITY_CONFIG_DIR"); custom != "" {
		return custom
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".antigravity-account-switcher"
	}
	return filepath.Join(home, ".config", "antigravity-account-switcher")
}

// ConfigFilePath returns the absolute path to the configuration file.
func ConfigFilePath() string {
	return filepath.Join(ConfigDir(), ConfigFileName)
}

// DefaultDBPath returns the default SQLite database path.
func DefaultDBPath() string {
	if env := os.Getenv("ANTIGRAVITY_DB_PATH"); env != "" {
		return env
	}
	return filepath.Join(ConfigDir(), DefaultDBFileName)
}

// DefaultConfig returns an initialized Config struct with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Port:                  DefaultPort,
		DBPath:                DefaultDBPath(),
		AntigravityBin:        "",
		UpstreamURL:           DefaultUpstreamURL,
		QuotaInterval:         DefaultInterval,
		QuotaWarningThreshold: DefaultQuotaWarningThreshold,
		QuotaSwitchThreshold:  DefaultQuotaSwitchThreshold,
		OpenBrowser:           false,
		ModelPrimary:          DefaultModelPrimary,
		ModelSecondary:        DefaultModelSecondary,
		FallbackSecondaryEnabled: DefaultFallbackSecondaryEnabled,
	}
}

// ParseBool parses string representations of boolean values.
func ParseBool(val string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(val)) {
	case "1", "t", "true", "yes", "y", "on":
		return true, nil
	case "0", "f", "false", "no", "n", "off":
		return false, nil
	default:
		return false, fmt.Errorf("cannot parse %q as boolean", val)
	}
}

// Load reads the configuration from disk, falling back to defaults if missing.
func Load() (*Config, error) {
	cfg := DefaultConfig()
	path := ConfigFilePath()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to read config file at %s: %w", path, err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON at %s: %w", path, err)
	}

	if cfg.QuotaWarningThreshold <= 0 {
		cfg.QuotaWarningThreshold = DefaultQuotaWarningThreshold
	}
	if cfg.QuotaSwitchThreshold <= 0 {
		cfg.QuotaSwitchThreshold = DefaultQuotaSwitchThreshold
	}

	// Environment variable overrides
	if envPort := os.Getenv("ANTIGRAVITY_PORT"); envPort != "" {
		var p int
		if _, err := fmt.Sscanf(envPort, "%d", &p); err == nil && p > 0 {
			cfg.Port = p
		}
	}
	if envDB := os.Getenv("ANTIGRAVITY_DB_PATH"); envDB != "" {
		cfg.DBPath = envDB
	}
	if envBin := os.Getenv("ANTIGRAVITY_BIN"); envBin != "" {
		cfg.AntigravityBin = envBin
	}
	if envUpstream := os.Getenv("ANTIGRAVITY_UPSTREAM_URL"); envUpstream != "" {
		cfg.UpstreamURL = envUpstream
	}
	if envPrimary := os.Getenv("ANTIGRAVITY_MODEL_PRIMARY"); envPrimary != "" {
		cfg.ModelPrimary = strings.TrimSpace(envPrimary)
	}
	if envSecondary := os.Getenv("ANTIGRAVITY_MODEL_SECONDARY"); envSecondary != "" {
		cfg.ModelSecondary = strings.TrimSpace(envSecondary)
	}
	if envFallback := os.Getenv("ANTIGRAVITY_FALLBACK_SECONDARY_ENABLED"); envFallback != "" {
		if b, err := ParseBool(envFallback); err == nil {
			cfg.FallbackSecondaryEnabled = b
		}
	}
	if envTunnelToken := os.Getenv("ANTIGRAVITY_CLOUDFLARE_TUNNEL_TOKEN"); envTunnelToken != "" {
		cfg.CloudflareTunnelToken = strings.TrimSpace(envTunnelToken)
	}
	if envAuthToken := os.Getenv("ANTIGRAVITY_REMOTE_AUTH_TOKEN"); envAuthToken != "" {
		cfg.RemoteAuthToken = strings.TrimSpace(envAuthToken)
	}

	if cfg.ModelPrimary == "" {
		cfg.ModelPrimary = DefaultModelPrimary
	}
	if cfg.ModelSecondary == "" {
		cfg.ModelSecondary = DefaultModelSecondary
	}

	return cfg, nil
}

// Save writes the configuration to disk, ensuring directory creation.
func Save(cfg *Config) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	path := ConfigFilePath()
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write config to %s: %w", path, err)
	}

	return nil
}

// CandidateAntigravityPaths returns standard probing locations across Linux and macOS.
func CandidateAntigravityPaths() []string {
	home, _ := os.UserHomeDir()
	candidates := []string{}

	if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
		candidates = append(candidates,
			filepath.Join(localAppData, "Programs", "Antigravity", "Antigravity.exe"),
			filepath.Join(localAppData, "Programs", "Antigravity", "antigravity.exe"),
			filepath.Join(localAppData, "Programs", "Antigravity", "bin", "antigravity.cmd"),
			filepath.Join(localAppData, "Programs", "Antigravity IDE", "Antigravity IDE.exe"),
			filepath.Join(localAppData, "Programs", "Antigravity IDE", "antigravity.exe"),
			filepath.Join(localAppData, "Programs", "Antigravity IDE", "bin", "antigravity.cmd"),
		)
	}
	if progFiles := os.Getenv("ProgramFiles"); progFiles != "" {
		candidates = append(candidates,
			filepath.Join(progFiles, "Antigravity", "Antigravity.exe"),
			filepath.Join(progFiles, "Antigravity", "antigravity.exe"),
			filepath.Join(progFiles, "Antigravity IDE", "Antigravity IDE.exe"),
			filepath.Join(progFiles, "Antigravity IDE", "antigravity.exe"),
		)
	}
	if progFilesX86 := os.Getenv("ProgramFiles(x86)"); progFilesX86 != "" {
		candidates = append(candidates,
			filepath.Join(progFilesX86, "Antigravity IDE", "antigravity.exe"),
			filepath.Join(progFilesX86, "Antigravity", "antigravity.exe"),
		)
	}

	if home != "" {
		candidates = append(candidates,
			// 1. Standard XDG user application directories (Recommended, no sudo)
			filepath.Join(home, ".local", "bin", "antigravity"),
			filepath.Join(home, ".local", "bin", "agy"),
			filepath.Join(home, ".local", "share", "antigravity", "antigravity"),
			filepath.Join(home, ".local", "share", "antigravity", "Antigravity-x64", "antigravity"),
			filepath.Join(home, ".local", "share", "Antigravity", "antigravity"),
			filepath.Join(home, ".local", "share", "Antigravity", "Antigravity-x64", "antigravity"),
			// 2. Custom tools & user development directories
			filepath.Join(home, "tools", "Antigravity", "Antigravity-x64", "antigravity"),
			filepath.Join(home, "tools", "Antigravity", "antigravity"),
			filepath.Join(home, "tools", "antigravity", "Antigravity-x64", "antigravity"),
			filepath.Join(home, "tools", "antigravity", "antigravity"),
			filepath.Join(home, "tools", "Antigravity-2.0", "antigravity"),
			filepath.Join(home, "Antigravity", "antigravity"),
			filepath.Join(home, "Antigravity", "Antigravity-x64", "antigravity"),
			filepath.Join(home, ".antigravity", "bin", "antigravity"),
			// 3. macOS Application Support
			filepath.Join(home, "Applications", "Antigravity.app", "Contents", "MacOS", "Antigravity"),
		)
	}

	// 4. System-wide FHS locations (Installed via sudo into /opt or /usr/local/bin)
	candidates = append(candidates,
		"/usr/local/bin/antigravity",
		"/usr/local/bin/agy",
		"/opt/antigravity/antigravity",
		"/opt/antigravity/Antigravity-x64/antigravity",
		"/opt/Antigravity/antigravity",
		"/opt/Antigravity/Antigravity-x64/antigravity",
		"/Applications/Antigravity.app/Contents/MacOS/Antigravity",
	)

	return candidates
}

// ResolveAntigravityBin determines the binary to run following priority:
// 1. Explicit override argument (if provided)
// 2. ANTIGRAVITY_BIN environment variable
// 3. Saved config file setting (`antigravity_bin`)
// 4. Candidate probing paths in filesystem
// 5. PATH lookups (`antigravity`, then `agy`)
func ResolveAntigravityBin(explicitOverride string) (string, error) {
	if explicitOverride != "" {
		if _, err := os.Stat(explicitOverride); err == nil {
			return filepath.Abs(explicitOverride)
		}
		// Check if it's a command in PATH
		if path, err := exec.LookPath(explicitOverride); err == nil {
			return path, nil
		}
		return "", fmt.Errorf("specified antigravity binary not found: %s", explicitOverride)
	}

	if env := os.Getenv("ANTIGRAVITY_BIN"); env != "" {
		if _, err := os.Stat(env); err == nil {
			return filepath.Abs(env)
		}
		if path, err := exec.LookPath(env); err == nil {
			return path, nil
		}
	}

	cfg, _ := Load()
	if cfg != nil && cfg.AntigravityBin != "" {
		if _, err := os.Stat(cfg.AntigravityBin); err == nil {
			return filepath.Abs(cfg.AntigravityBin)
		}
		if path, err := exec.LookPath(cfg.AntigravityBin); err == nil {
			return path, nil
		}
	}

	// Probe candidate filesystem locations
	for _, candidate := range CandidateAntigravityPaths() {
		if fi, err := os.Stat(candidate); err == nil && !fi.IsDir() {
			// On Windows, os.Stat does not report POSIX exec bits for .exe files.
			if runtime.GOOS == "windows" || fi.Mode()&0o111 != 0 {
				return filepath.Abs(candidate)
			}
		}
	}

	// Look up in PATH
	for _, cmdName := range []string{"antigravity", "agy"} {
		if path, err := exec.LookPath(cmdName); err == nil {
			return path, nil
		}
	}

	return "", errors.New("could not automatically locate Antigravity binary; specify it via 'antigravity-account-switcher config set antigravity_bin <path>' or flag '--bin'")
}

// FindAntigravityIcon searches for an icon in the directory of the resolved binary.
func FindAntigravityIcon(resolvedBin string) string {
	if resolvedBin == "" {
		return ""
	}

	dir := filepath.Dir(resolvedBin)
	searchDirs := []string{
		dir,
		filepath.Dir(dir), // Parent dir (e.g. ~/tools/Antigravity if bin is in Antigravity-x64)
		filepath.Join(dir, "resources"),
	}

	for _, d := range searchDirs {
		for _, name := range []string{"icon.png", "logo.png", "antigravity.png"} {
			p := filepath.Join(d, name)
			if fi, err := os.Stat(p); err == nil && !fi.IsDir() && fi.Size() > 0 {
				return p
			}
		}
	}

	return ""
}
