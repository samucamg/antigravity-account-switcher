# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-09-04

### Added
- **Proactive Quota Switching (85% Threshold):** Reverse proxy automatically rotates active accounts when model quota usage reaches 85%, eliminating `HTTP 429 RESOURCE_EXHAUSTED` interruptions before they occur.
- **Real-Time Quota Warning Alerts (80% Threshold):** Quota poller and event bus emit `EventTypeQuotaWarning` when accounts reach 80% usage, broadcasted via SSE to the Web Dashboard.
- **Configurable Thresholds:** Added `quota_warning_threshold` (default `0.80`) and `quota_switch_threshold` (default `0.85`) to `config.json` and CLI configuration commands.
- **Unit Tests:** Coverage for `QuotaBucket.UsageFraction`, `QuotaBucket.IsUsageAboveThreshold`, `FailoverEngine.RotateProactively`, and quota warning event emission.
- **Windows OS Native Support:** Comprehensive documentation and path detection for native Windows (`.exe`) execution alongside existing Linux (XDG/FHS) and WSL2 workflows.
- **Language Flag Badges:** Added Brazilian Portuguese (🇧🇷) and English (🇺🇸) navigation flags to documentation.
- **Beautified Documentation:** Reworked READMEs (EN + PT-BR) with platform badges, feature overview table, quick-start, and friendly emoji-rich formatting.

### Changed
- **Bidirectional Tunnel Synchronization (`handleConnect`):** Refactored HTTP CONNECT tunnel to use `sync.WaitGroup` and `CloseWrite()`, preventing premature socket termination during half-closed transfers.
- **Lock Context Awareness (`RotateAccount`):** Added immediate `ctx.Err()` check upon acquiring mutex in `FailoverEngine`, avoiding unnecessary database rotations for aborted client requests.
- **System Proxy Compliance:** Configured `Proxy: http.ProxyFromEnvironment` and standard dialers across `ProxyHandler` and `QuotaPoller` to respect `HTTP_PROXY`, `HTTPS_PROXY`, and `NO_PROXY` environment variables.

### Security
- **Audited Codebase:** Verified zero telemetry exfiltration, strict token redaction (`json:"-"`), and secure parameterization across all SQLite queries.
