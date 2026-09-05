# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2026-09-05

### Added
- **Multi-Tier Model Fallback:** Zero-allocation in-flight request rewriter and intelligent intra-account failover engine. When primary model quota exhausts (HTTP 429 / 403 RESOURCE_EXHAUSTED), requests seamlessly fall back to secondary model (e.g., Gemini Pro -> Flash) on the same account before rotating accounts.
- **Model Routing Dashboard & CLI:** Web dashboard controls with toggle, primary/secondary model selectors, live model discovery from local Language Server/Cloud Code, and CLI flags (`--fallback-secondary`, `--model-primary`, `--model-secondary`).
- **Cloudflare Tunnels Integration:** Hybrid tunneling module supporting 1-click Quick Tunnels (`trycloudflare.com`) without account requirements, as well as Cloudflare Zero Trust Named Tunnels with tokens for custom fixed domains.
- **1-Click Windows Launcher (`start.bat`):** Standalone zero-setup batch launcher with dynamic IDE path resolution (no hardcoded user paths), silent background daemon, and automated browser launch.
- **Automated Desktop Shortcut:** `install.ps1` automatically creates `Antigravity (Multi-Account).lnk` on the Windows Desktop with the official Antigravity icon.
- **Client Timezone Awareness:** SQLite strftime daily bucket grouping now respects client local timezone (`?tz=...` and `tz_offset=...`), ensuring token usage metrics accurately reflect user local calendar days.
- **AI Agent Harmonization:** Added canonical `AGENTS.md` spec and IDE rules for Cursor (`.cursorrules`), Windsurf (`.windsurfrules`), Claude Code (`CLAUDE.md`), and GitHub Copilot (`copilot-instructions.md`).

### Changed
- **Extended Quota Network Resilience:** Increased quota client HTTP timeout from 15s to 30s (`DefaultHTTPTimeout`) to absorb upstream Google Cloud Code PA latency spikes, and adjusted default polling interval to 5m.
- **Format Target:** Added `make fmt` target (`gofmt -s -w .`) to Makefile.

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
