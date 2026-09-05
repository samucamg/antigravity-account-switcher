<p align="center">
  <img src="https://img.shields.io/badge/Made%20with-Go-00ADD8?style=plastic&logo=go&logoColor=white" alt="Go"/>
  <img src="https://img.shields.io/badge/version-v0.2.0-blue?style=plastic" alt="version"/>
  <img src="https://img.shields.io/badge/status-stable-success?style=plastic" alt="status"/>
</p>

<h1 align="center">🚀 Antigravity Account Switcher</h1>

<p align="center"><b>Your “never get rate-limited again” sidekick for Google Antigravity 2.0.</b></p>

<p align="center">
  🇺🇸 English · <a href="README.pt-BR.md">🇧🇷 Português (Brasil)</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Platform-Windows_%7C_Linux_%7C_macOS-6f42c1?style=plastic&logo=windows&logoColor=white" alt="platform"/>
  <img src="https://img.shields.io/github/go-mod/go-version/samucamg/antigravity-account-switcher?style=plastic&logo=go&label=Go" alt="go version"/>
  <img src="https://img.shields.io/github/license/samucamg/antigravity-account-switcher?style=plastic&color=blue" alt="license"/>
  <img src="https://img.shields.io/github/stars/samucamg/antigravity-account-switcher?style=plastic&logo=github" alt="stars"/>
  <img src="https://github.com/samucamg/antigravity-account-switcher/actions/workflows/ci.yml/badge.svg" alt="CI"/>
</p>

<p align="center">
  🧠 Claude &nbsp;·&nbsp; 🤖 GPT &nbsp;·&nbsp; ✨ Gemini — <b>all behind an intelligent multi-account pool</b>
</p>

---

## 😩 Has this ever happened to you?

You are in the middle of an awesome coding session with Google Antigravity 2.0… and suddenly:

```text
HTTP 429 RESOURCE_EXHAUSTED
```

Quota exhausted. Flow interrupted. Momentum lost. 😤

**Antigravity Account Switcher** exists so this never happens to you again.

---

## ✨ What does it do for you?

The switcher is a **transparent, zero-latency local supervisor** that sits between Antigravity 2.0 and Google Cloud Code:

| 💜 Feature | What happens |
|:---|:---|
| 🔔 **Warning at 80%** | When an account hits 80% quota usage, you get an instant real-time alert (web dashboard + events) — enough time to plan calmly. |
| 🔁 **Proactive switch at 85%** | Before quota exhaustion disrupts your workflow, the switcher seamlessly rotates to the next healthy account in your pool. |
| 🧠 **Multi-Tier Model Fallback** *(v0.2.0)* | Exhausted a premium model (e.g. Gemini 2.5 Pro)? The proxy dynamically rewrites in-flight requests to a secondary model (e.g. Gemini 2.5 Flash) on the **same** account before rotating accounts! |
| 🔄 **Reactive failover (HTTP 429)** | If a 429 still slips through, the in-flight request is replayed in memory with the new account — zero editor error dialogs. |
| 🚇 **Cloudflare Tunnels** *(v0.2.0)* | Built-in 1-click tunneling: create instant public Quick Tunnels (`trycloudflare.com`) or connect custom Zero Trust named tunnels with tokens for secure remote access. |
| ⚡ **1-Click Windows Launcher** *(v0.2.0)* | Standalone `start.bat` and desktop shortcut automatically start the proxy daemon, open the web dashboard, and launch Antigravity IDE. |
| 📊 **Real-time Quotas & Timezone Metrics** | Live model quotas (Claude, GPT, Gemini) and request history charts rendered accurately in your local timezone. |
| 🔒 **100% local & private** | No token or prompt ever leaves your machine. Everything lives securely in your local SQLite database (`accounts.db`). |
| 🪟🐧🍏 **Windows, Linux & macOS** | High-performance, single dependency-free Go binary. |

---

## 🛠️ How it works (in 10 seconds)

```text
  Google Antigravity 2.0 / agy
          │  (requests intercepted by local proxy on port 1831 / 8080)
          ▼
  🧲 Antigravity Account Switcher
          │  uses active account from pool
          │
          ├─ quota ≥ 80%             → 🔔 fires warning alert
          ├─ quota ≥ 85%             → 🔁 proactive account rotation
          ├─ model quota exhausted   → 🧠 fallback to secondary model (Pro → Flash)
          └─ HTTP 429 response       → 🔄 rotate account + in-memory replay
          ▼
  ☁️ Google Cloud Code (Cloud Code PA)
```

---

## 📦 Installation Guide

### 🪟 Windows (Recommended & Quickest Setup)

We provide an automated installer and 1-click launcher designed specifically for Windows users.

#### 1. Prerequisites
- **Go 1.24+**: Download the Windows MSI installer from [go.dev/dl](https://go.dev/dl/) and install it. Verify by opening PowerShell and typing `go version`.
- **Git for Windows**: Download from [git-scm.com](https://git-scm.com/) (or run `winget install Git.Git`).
- **Google Antigravity 2.0 IDE**: Installed on your system (standard path is `%LOCALAPPDATA%\Programs\Antigravity\Antigravity.exe`).
- *(Optional - for Cloudflare Tunnels)*: If you want remote access tunneling, install `cloudflared` via PowerShell:
  ```powershell
  winget install --id Cloudflare.cloudflared
  ```

---

#### 2. Automated 1-Step Installation (`install.ps1`)

Open PowerShell, navigate to the folder where you want to keep the tool, and run:

```powershell
git clone https://github.com/samucamg/antigravity-account-switcher.git
cd antigravity-account-switcher
powershell -ExecutionPolicy Bypass -File .\install.ps1
```

> [!TIP]
> If PowerShell gives an error saying *“running scripts is disabled on this system”*, you can allow local script execution by running:
> ```powershell
> Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
> ```
> And then re-run `.\install.ps1`.

**What `install.ps1` does automatically for you:**
1. ✅ Checks and validates your Go compiler.
2. 🔨 Compiles `antigravity-account-switcher.exe` with optimizations.
3. ⚙️ Configures the default service port to `1831`.
4. 🔍 Automatically discovers your Google Antigravity 2.0 IDE executable path.
5. 📄 Creates helper scripts: `start.ps1`, `add-account.ps1`, `switch-account.ps1`.
6. 🖥️ **Creates a Desktop Shortcut:** Adds an **`Antigravity (Multi-Account).lnk`** shortcut right on your Windows Desktop, with the official Antigravity IDE icon!

---

#### 3. Daily Usage on Windows (1-Click)

You have two convenient ways to start your work session:

- **Method A (Easiest):** Simply double-click the **`Antigravity (Multi-Account)`** shortcut on your Desktop!
- **Method B:** Double-click **`start.bat`** inside the `antigravity-account-switcher` directory.

**What happens behind the scenes:**
- 🛡️ Closes any previous stale switcher processes to prevent port conflicts.
- 🚀 Launches the switcher proxy invisibly in the background.
- 🌐 Opens your default web browser directly to the Web Dashboard at `http://127.0.0.1:1831/`.
- 🔗 Automatically injects the proxy configuration (`HTTP_PROXY`, `HTTPS_PROXY`, `CLOUD_CODE_URL`) and launches **Google Antigravity 2.0 IDE**!

---

#### 4. Manual CLI Usage on Windows (PowerShell / CMD)

If you prefer using the terminal directly:

```powershell
# Compile the binary
go build -o antigravity-account-switcher.exe ./cmd/antigravity-account-switcher

# Add a Google account (opens browser for secure OAuth2)
.\antigravity-account-switcher.exe add-account

# Launch Antigravity IDE under the supervised proxy
.\antigravity-account-switcher.exe launch

# Or run only the background proxy and web dashboard
.\antigravity-account-switcher.exe serve --port 1831
```

---

#### 5. Windows Troubleshooting & FAQ

- **Port 1831 already in use?**
  `start.bat` automatically terminates leftover processes. If another application occupies port 1831, change the port anytime:
  ```powershell
  .\antigravity-account-switcher.exe config set port 1832
  ```
- **Antigravity IDE installed in a custom directory?**
  If `install.ps1` didn't find your IDE automatically, set your custom path:
  ```powershell
  .\antigravity-account-switcher.exe config set antigravity_bin "D:\Tools\Antigravity\Antigravity.exe"
  ```
- **Windows Defender Firewall prompt:**
  On first run, Windows may prompt you about network access for `antigravity-account-switcher.exe`. Select **Private networks** and click **Allow access** so that Antigravity IDE can communicate with the local proxy.
- **Using WSL2 (Windows Subsystem for Linux)?**
  You can run the switcher completely inside WSL2! Follow the Linux steps below. The switcher automatically detects and calls the Windows Antigravity binary across the WSL boundary.

---

### 🐧 Linux / 🍏 macOS

```bash
# 1. Clone repository
git clone https://github.com/samucamg/antigravity-account-switcher.git
cd antigravity-account-switcher

# 2. Build and install binary (installs to ~/.local/bin)
make install

# 3. (Optional on Linux) Create desktop app menu entry
antigravity-account-switcher install-desktop
```

---

## 🏃 Quick Start in 3 Steps

1️⃣ **Add your Google accounts** (as many as you want — 2, 3, 10…):

```bash
# Windows
.\antigravity-account-switcher.exe add-account

# Linux / macOS
antigravity-account-switcher add-account
```
*(Or click **“Add Account”** directly inside the Web Dashboard!)*

2️⃣ **Start working with Antigravity under supervision:**
- On **Windows**: Double-click the **Antigravity (Multi-Account)** Desktop shortcut or `start.bat`.
- On **Linux / macOS**: Run `antigravity-account-switcher launch`.

3️⃣ **Monitor live quotas & metrics:**
Open `http://127.0.0.1:1831/` (or `http://127.0.0.1:8080/`) in your browser to view accounts, quota gauges, token charts, and switch logs in real time. 👀

> 💡 **Auto-Import:** Already logged into Google Antigravity? On first run, the switcher **automatically imports your active session credentials** — zero manual setup required.

---

## 🌟 New Features in v0.2.0

### 🧠 Multi-Tier Model Fallback
When working on heavy coding tasks, high-end models (such as `gemini-2.5-pro` or Claude 3.5 Sonnet) may deplete their rolling quota while lighter models on the *same* account still have 100% capacity.
- **Intelligent intra-account failover:** When primary model quota runs out (HTTP 429 / 403 `RESOURCE_EXHAUSTED`), the switcher seamlessly rewrites the in-flight request to a secondary fallback model (e.g. `gemini-2.5-flash`) on the **same account**, without prematurely cycling to your next account!
- **Zero-allocation streaming:** Uses an optimized streaming rewriter that modifies model headers and JSON payload bodies in flight with negligible CPU and memory overhead.
- **Configurable via Web UI or CLI:**
  ```bash
  antigravity-account-switcher serve --fallback-secondary --model-primary gemini-2.5-pro --model-secondary gemini-2.5-flash
  ```

---

### 🚇 Cloudflare Tunnels (Remote Access & Webhooks)
Need to access your Antigravity Account Switcher dashboard remotely from your phone, tablet, or another workstation?
- **Quick Tunnels (1-Click):** Start an instant, zero-configuration HTTPS tunnel via `trycloudflare.com`. No Cloudflare account or credit card required.
- **Zero Trust Named Tunnels:** Connect a permanent custom domain via Cloudflare Zero Trust tunnel token (`eyJh...`) with access policies and authentication.
- **Web Dashboard Integration:** Toggle tunnels on/off with 1 click, view live connection status, and copy public URLs instantly.
- **Prerequisite:** Ensure `cloudflared` is installed (`winget install --id Cloudflare.cloudflared` on Windows, or `brew install cloudflared` on macOS).

---

## 🖥️ Web Dashboard

Access the interactive dashboard at `http://127.0.0.1:1831/`:

- 📋 **Account Pool Overview:** Status of every account (Active, Available, Exhausted, Cooldown).
- 🧮 **Per-Model Quotas:** Real-time percentage bars for Gemini, Claude, and GPT models.
- 🧠 **Model Fallback Routing:** Toggle intra-account model fallback and pick primary/secondary models with live discovery.
- 🚇 **Cloudflare Tunnel Control:** 1-click Quick Tunnel generation and Zero Trust tunnel management.
- 🔔 **Visual Alerts:** Visual warning cues when usage hits 80%.
- 📈 **Timezone-Aware Charts:** Hourly and daily token consumption graphs mapped directly to your local clock.
- ➕ **Web OAuth2 Onboarding:** Add new Google accounts directly from your browser with a single button click.

---

## ⌨️ CLI Command Reference

| Command | Description |
|:---|:---|
| `launch` | **(Recommended)** Launches Antigravity 2.0 under the supervised proxy |
| `serve` | Runs the proxy engine + quota monitor + web dashboard |
| `add-account` | Adds a new Google account via OAuth2 (loopback + PKCE) |
| `list-accounts` | Lists all accounts, the active one, and remaining quota % |
| `switch-account [email]` | Manually switches the active account in the pool |
| `refresh-quotas` | Forces an immediate quota synchronization with Google |
| `wrap -- <cmd>` | Executes any arbitrary command with proxy environment variables |
| `status` | Displays active account, token health, and supervisor metrics |
| `config` | View and modify configuration values (`get`, `set`, `list`) |
| `install-desktop` | Creates application menu icon on Linux desktops |
| `uninstall-desktop` | Removes the Linux application menu icon |
| `version` | Displays version (`v0.2.0`), commit hash, and build timestamp |

### CLI Flags for `serve` and `launch`
- `--port <port>`: Port to listen on (e.g., `1831` or `8080`).
- `--fallback-secondary`: Enables intra-account secondary model fallback before account rotation.
- `--model-primary <id>`: Primary model ID (default: `gemini-2.5-pro` or auto-detected).
- `--model-secondary <id>`: Fallback model ID (default: `gemini-2.5-flash`).

---

## ⚙️ Configuration

Configuration is stored in:
- **Windows:** `%APPDATA%\antigravity-account-switcher\config.json`
- **Linux / macOS:** `~/.config/antigravity-account-switcher/config.json`

```bash
# View all settings
antigravity-account-switcher config list

# Customize settings
antigravity-account-switcher config set port 1831
antigravity-account-switcher config set quota_interval 60s
antigravity-account-switcher config set antigravity_bin "C:\Path\To\Antigravity.exe"

# Quota thresholds (0.0 to 1.0)
antigravity-account-switcher config set quota_warning_threshold 0.80   # 🔔 warn at 80%
antigravity-account-switcher config set quota_switch_threshold  0.85   # 🔁 switch at 85%
```

**Environment Variables:** `ANTIGRAVITY_BIN`, `ANTIGRAVITY_PORT`, `ANTIGRAVITY_DB_PATH`, `ANTIGRAVITY_CLIENT_ID`, `ANTIGRAVITY_CLIENT_SECRET`.

---

## 🔒 Security & Privacy

- 🛡️ **IETF RFC Standard OAuth2:** Secure loopback flow with PKCE (RFC 7636) — the safest authorization flow for native apps.
- 🏠 **100% Local Storage:** Credentials, refresh tokens, and metadata are saved strictly inside your local SQLite database (`accounts.db`) with restricted file permissions.
- 🚫 **Zero Telemetry:** No prompts, tokens, metrics, or personal data leave your computer.
- 🔐 **Transparent Tunneling:** Non-Cloud Code traffic (such as voice STT to `speech.googleapis.com`) is tunneled via raw bidirectional TCP `CONNECT` without inspection.

See [SECURITY.md](SECURITY.md) for more details.

---

## 🤝 Contributing

Contributions are welcome! Please check out [CONTRIBUTING.md](CONTRIBUTING.md) to set up your development environment, run race-detector tests, and submit a pull request.

```bash
# Run tests with race detection
go test -race ./...
```

---

## ⭐ Enjoying it?

Give this repo a ⭐ on GitHub and share it with fellow developers who are tired of getting rate-limited! 😄

---

## 📄 License

MIT © 2026 **Muriel Gasparini** — see [LICENSE](LICENSE).
