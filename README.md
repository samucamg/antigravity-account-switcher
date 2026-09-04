<p align="center">
  <img src="https://img.shields.io/badge/Made%20with-Go-00ADD8?style=plastic&logo=go&logoColor=white" alt="Go"/>
  <img src="https://img.shields.io/badge/status-stable-success?style=plastic" alt="status"/>
</p>

<h1 align="center">🚀 Antigravity Account Switcher</h1>

<p align="center"><b>Your “never get rate-limited again” sidekick for Google Antigravity 2.0.</b></p>

<p align="center">
  🇺🇸 English · <a href="README.pt-BR.md">🇧🇷 Português (Brasil)</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Platform-Linux_%7C_Windows_%7C_macOS-6f42c1?style=plastic&logo=linux&logoColor=white" alt="platform"/>
  <img src="https://img.shields.io/github/go-mod/go-version/Muriel-Gasparini/antigravity-account-switcher?style=plastic&logo=go&label=Go" alt="go version"/>
  <img src="https://img.shields.io/github/license/Muriel-Gasparini/antigravity-account-switcher?style=plastic&color=blue" alt="license"/>
  <img src="https://img.shields.io/github/stars/Muriel-Gasparini/antigravity-account-switcher?style=plastic&logo=github" alt="stars"/>
  <img src="https://github.com/Muriel-Gasparini/antigravity-account-switcher/actions/workflows/ci.yml/badge.svg" alt="CI"/>
</p>

<p align="center">
  🧠 Claude &nbsp;·&nbsp; 🤖 GPT &nbsp;·&nbsp; ✨ Gemini — <b>all behind a single account pool</b>
</p>

---

## 😩 Has this ever happened to you?

You are in the middle of an awesome coding session with Antigravity 2.0… and suddenly:

```text
HTTP 429 RESOURCE_EXHAUSTED
```

Quota exhausted. Flow interrupted. Momentum lost. 😤

**Antigravity Account Switcher** exists so this never happens to you again.

---

## ✨ What does it do for you?

The switcher is a **transparent local supervisor** that sits between Antigravity 2.0 and Google Cloud Code:

| 💜 Feature | What happens |
|:---|:---|
| 🔔 **Warning at 80%** | When an account hits 80% usage you get a real-time alert (web dashboard + events) — enough time to plan calmly. |
| 🔁 **Proactive switch at 85%** | Before the quota blows up, the switcher rotates to the next healthy account in your pool. Your session never notices. |
| 🔄 **Reactive failover (HTTP 429)** | If a 429 still slips through, the in-flight request is replayed in memory with the new account — no editor errors. |
| 📊 **Real-time quotas** | Monitors usage per model (Claude, GPT, Gemini) and renders everything in a beautiful web dashboard. |
| 🔒 **100% local & private** | No token or prompt ever leaves your machine. Everything lives in your local SQLite (`accounts.db`). |
| 🪟🐧🍏 **Linux, Windows & macOS** | A single dependency-free Go binary. |

---

## 🛠️ How it works (in 10 seconds)

```text
  Antigravity 2.0 / agy
          │  (request intercepted by the local proxy)
          ▼
  🧲 Antigravity Account Switcher
          │  uses the active account from the pool
          │
          ├─ quota ≥ 80%  → 🔔 fires warning alert
          ├─ quota ≥ 85%  → 🔁 rotates account (on the spot)
          └─ HTTP 429     → 🔄 tries next account + replays in memory
          ▼
  ☁️ Google Cloud Code (Cloud Code PA)
```

---

## 📦 Installing the Switcher

### 🐧 Linux / macOS

```bash
git clone https://github.com/Muriel-Gasparini/antigravity-account-switcher.git
cd antigravity-account-switcher
make install        # static binary → ~/.local/bin
```

### 🪟 Windows

**Option A — build on Windows (PowerShell or CMD):**

```powershell
git clone https://github.com/Muriel-Gasparini/antigravity-account-switcher.git
cd antigravity-account-switcher
go build -o antigravity-account-switcher.exe ./cmd/antigravity-account-switcher
.\antigravity-account-switcher.exe serve
```

**Option B — using WSL2 (Windows Subsystem for Linux):**
just follow the Linux steps inside WSL2. The switcher auto-detects the Antigravity binary installed on Windows.

> [!TIP]
> On Windows, Antigravity is usually installed at
> `%LOCALAPPDATA%\Programs\Antigravity\antigravity.exe` — and the switcher finds that path automatically. 🎯

---

## 🏃 Quick start in 3 steps

1️⃣ **Add your Google accounts** (as many as you want — 2, 3, 10…)

```bash
antigravity-account-switcher add-account
```

2️⃣ **Launch Antigravity under supervision** (recommended) or run the service only:

```bash
antigravity-account-switcher launch   # proxy + opens Antigravity
antigravity-account-switcher serve    # proxy + dashboard only
```

3️⃣ **Keep an eye on the dashboard:** open `http://127.0.0.1:8080/` in your browser and watch accounts, quotas and history live. 👀

> 💡 Already logged into Antigravity? On first run the switcher **imports your login automatically** — zero configuration.

---

## 🖥️ Web Dashboard

- 📋 status of every account in the pool (active, available, exhausted)
- 🧮 remaining quota percentage per model
- 🔔 visual alerts when usage hits 80%
- 📈 token charts and switch history (real-time SSE)

---

## ⌨️ CLI commands

| Command | Description |
|:---|:---|
| `launch` | **(recommended)** starts Antigravity 2.0 under the supervised proxy |
| `serve` | runs proxy + quota monitor + web dashboard |
| `wrap -- <cmd>` | runs any command with the switcher proxy |
| `add-account` | adds a Google account via OAuth2 (loopback + PKCE) |
| `list-accounts` | lists accounts, active one and % quota |
| `refresh-quotas` | syncs quotas with Google on demand |
| `status` | shows active account, tokens and service health |
| `config` | manages settings (`get`, `set`, `list`) |
| `install-desktop` | creates a launcher icon entry on Linux (GNOME/KDE/XFCE) |
| `uninstall-desktop` | removes the launcher entry |
| `version` | binary version, commit and build date |

---

## ⚙️ Configuration

File: `~/.config/antigravity-account-switcher/config.json` — or `%APPDATA%\antigravity-account-switcher\config.json` on Windows.

```bash
antigravity-account-switcher config list

antigravity-account-switcher config set port 8080
antigravity-account-switcher config set quota_interval 60s
antigravity-account-switcher config set antigravity_bin /path/to/antigravity

# quota thresholds (values from 0.0 to 1.0)
antigravity-account-switcher config set quota_warning_threshold 0.80   # 🔔 warn at 80%
antigravity-account-switcher config set quota_switch_threshold  0.85   # 🔁 switch at 85%
```

**Environment variables:** `ANTIGRAVITY_BIN`, `ANTIGRAVITY_PORT`, `ANTIGRAVITY_DB_PATH`, `ANTIGRAVITY_CLIENT_ID`, `ANTIGRAVITY_CLIENT_SECRET`.

---

## 🔒 Security & Privacy

- 🛡️ **Secure OAuth2:** RFC 8252 loopback flow with PKCE (RFC 7636) — the IETF standard for native apps.
- 🏠 **100% local:** credentials and tokens live only inside your `accounts.db` (SQLite) with restricted permissions.
- 🚫 **Zero telemetry:** no data leaves your machine to third parties.

See [SECURITY.md](SECURITY.md) for more.

---

## 🤝 Contributing

We love contributions! 🎉 Check [CONTRIBUTING.md](CONTRIBUTING.md) to set up the environment, run the race-detector tests and open your PR.

```bash
go test ./...
```

---

## ⭐ Enjoying it?

Drop a ⭐ on the repo and share it with everyone who keeps getting 429s! 😄

---

## 📄 License

MIT © 2026 **Muriel Gasparini** — see [LICENSE](LICENSE).
