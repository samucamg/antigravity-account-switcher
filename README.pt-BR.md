<p align="center">
  <img src="https://img.shields.io/badge/Feito%20com-Go-00ADD8?style=plastic&logo=go&logoColor=white" alt="Go"/>
  <img src="https://img.shields.io/badge/status-est%C3%A1vel-success?style=plastic" alt="status"/>
</p>

<h1 align="center">🚀 Antigravity Account Switcher</h1>

<p align="center"><b>Seu “nunca mais fique sem cota” para o Google Antigravity 2.0.</b></p>

<p align="center">
  <a href="README.md">🇺🇸 English</a> · 🇧🇷 Português (Brasil)
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Plataforma-Linux_%7C_Windows_%7C_macOS-6f42c1?style=plastic&logo=linux&logoColor=white" alt="platform"/>
  <img src="https://img.shields.io/github/go-mod/go-version/Muriel-Gasparini/antigravity-account-switcher?style=plastic&logo=go&label=Go" alt="go version"/>
  <img src="https://img.shields.io/github/license/Muriel-Gasparini/antigravity-account-switcher?style=plastic&color=blue" alt="license"/>
  <img src="https://img.shields.io/github/stars/Muriel-Gasparini/antigravity-account-switcher?style=plastic&logo=github" alt="stars"/>
  <img src="https://github.com/Muriel-Gasparini/antigravity-account-switcher/actions/workflows/ci.yml/badge.svg" alt="CI"/>
</p>

<p align="center">
  🧠 Claude &nbsp;·&nbsp; 🤖 GPT &nbsp;·&nbsp; ✨ Gemini — <b>todos sob uma única fila de contas</b>
</p>

---

## 😩 Você já passou por isso?

Você está no meio de uma sessão incrível de programação no Antigravity 2.0… e, do nada:

```text
HTTP 429 RESOURCE_EXHAUSTED
```

Cota esgotada. Fluxo interrompido. Raciocínio perdido. 😤

O **Antigravity Account Switcher** existe para que isso nunca mais aconteça com você.

---

## ✨ O que ele faz por você?

O switcher é um **supervisor local transparente** que fica entre o Antigravity 2.0 e o Google Cloud Code:

| 💜 Recurso | O que acontece |
|:---|:---|
| 🔔 **Alerta aos 80%** | Quando uma conta atinge 80% do uso, você recebe um aviso em tempo real (dashboard web + eventos) — dá tempo de agir com calma. |
| 🔁 **Troca proativa aos 85%** | Antes de estourar a cota, o switcher rotaciona para a próxima conta saudável do seu pool. A sessão nem percebe. |
| 🔄 **Failover reativo (HTTP 429)** | Se mesmo assim vier um 429, a requisição em voo é repetida em memória com a nova conta — sem erro no editor. |
| 📊 **Cotas em tempo real** | Monitora o consumo por modelo (Claude, GPT, Gemini) e exibe tudo numa dashboard web linda. |
| 🔒 **100% local e privado** | Nenhum token ou prompt sai da sua máquina. Tudo fica no seu SQLite local (`accounts.db`). |
| 🪟🐧🍏 **Linux, Windows e macOS** | Um único binário Go, sem dependências externas. |

---

## 🛠️ Como funciona (em 10 segundos)

```text
  Antigravity 2.0 / agy
          │  (pedido interceptado pelo proxy local)
          ▼
  🧲 Antigravity Account Switcher
          │  usa a conta ativa do pool
          │
          ├─ cota ≥ 80%  → 🔔 dispara alerta
          ├─ cota ≥ 85%  → 🔁 troca de conta (na hora)
          └─ HTTP 429    → 🔄 tenta próxima conta + repete em memória
          ▼
  ☁️ Google Cloud Code (Cloud Code PA)
```

---

## 📦 Instalação do Switcher

### 🐧 Linux / macOS

```bash
git clone https://github.com/Muriel-Gasparini/antigravity-account-switcher.git
cd antigravity-account-switcher
make install        # binário estático → ~/.local/bin
```

### 🪟 Windows

**Opção A — compilar no Windows (PowerShell ou CMD):**

```powershell
git clone https://github.com/Muriel-Gasparini/antigravity-account-switcher.git
cd antigravity-account-switcher
go build -o antigravity-account-switcher.exe ./cmd/antigravity-account-switcher
.\antigravity-account-switcher.exe serve
```

**Opção B — usando WSL2 (Windows Subsystem for Linux):**
siga os passos de Linux normalmente dentro do WSL2. O switcher detecta sozinho o caminho do Antigravity instalado no Windows.

> [!TIP]
> No Windows, o Antigravity normalmente instala em
> `%LOCALAPPDATA%\Programs\Antigravity\antigravity.exe` — e o switcher encontra esse caminho automaticamente. 🎯

---

## 🏃 Comece em 3 passos

1️⃣ **Adicione suas contas Google** (quantas quiser — 2, 3, 10…)

```bash
antigravity-account-switcher add-account
```

2️⃣ **Lance o Antigravity sob supervisão** (recomendado) ou rode só o serviço:

```bash
antigravity-account-switcher launch   # proxy + abre o Antigravity
antigravity-account-switcher serve    # só o proxy + dashboard
```

3️⃣ **Acompanhe tudo pela dashboard:** abra `http://127.0.0.1:8080/` no navegador e veja contas, cotas e histórico em tempo real. 👀

> 💡 Já tem o Antigravity logado? Na primeira execução o switcher **importa seu login automaticamente** — zero configuração.

---

## 🖥️ Dashboard Web

- 📋 status de cada conta do pool (ativa, disponível, exausta)
- 🧮 percentual de cota restante por modelo
- 🔔 alertas visuais ao atingir 80% de uso
- 📈 gráficos de tokens e histórico de trocas (SSE em tempo real)

---

## ⌨️ Comandos da CLI

| Comando | Descrição |
|:---|:---|
| `launch` | **(recomendado)** inicia o Antigravity 2.0 sob o proxy supervisionado |
| `serve` | roda proxy + monitor de cotas + dashboard web |
| `wrap -- <cmd>` | executa qualquer comando com o proxy do switcher |
| `add-account` | adiciona uma conta Google via OAuth2 (loopback + PKCE) |
| `list-accounts` | lista contas, ativa e % de cota |
| `refresh-quotas` | sincroniza cotas na hora com o Google |
| `status` | mostra conta ativa, tokens e saúde do serviço |
| `config` | gerencia configurações (`get`, `set`, `list`) |
| `install-desktop` | cria atalho com ícone no menu do Linux (GNOME/KDE/XFCE) |
| `uninstall-desktop` | remove o atalho |
| `version` | versão, commit e data do binário |

---

## ⚙️ Configuração

Arquivo: `~/.config/antigravity-account-switcher/config.json` — ou `%APPDATA%\antigravity-account-switcher\config.json` no Windows.

```bash
antigravity-account-switcher config list

antigravity-account-switcher config set port 8080
antigravity-account-switcher config set quota_interval 60s
antigravity-account-switcher config set antigravity_bin /caminho/do/antigravity

# limiares de cota (valores de 0.0 a 1.0)
antigravity-account-switcher config set quota_warning_threshold 0.80   # 🔔 alerta aos 80%
antigravity-account-switcher config set quota_switch_threshold  0.85   # 🔁 troca aos 85%
```

**Variáveis de ambiente:** `ANTIGRAVITY_BIN`, `ANTIGRAVITY_PORT`, `ANTIGRAVITY_DB_PATH`, `ANTIGRAVITY_CLIENT_ID`, `ANTIGRAVITY_CLIENT_SECRET`.

---

## 🔒 Segurança & Privacidade

- 🛡️ **OAuth2 seguro:** fluxo loopback RFC 8252 com PKCE (RFC 7636) — padrão IETF para apps nativos.
- 🏠 **100% local:** credenciais e tokens vivem apenas no seu `accounts.db` (SQLite) com permissões restritas.
- 🚫 **Zero telemetria:** nada de dados sai da sua máquina para terceiros.

Veja mais em [SECURITY.md](SECURITY.md).

---

## 🤝 Contribuindo

Adoramos receber ajuda! 🎉 Veja [CONTRIBUTING.md](CONTRIBUTING.md) para configurar o ambiente, rodar os testes com detector de *race* e abrir seu PR.

```bash
go test ./...
```

---

## ⭐ Gostou?

Deixa uma ⭐ no repositório e compartilhe com quem vive tomando 429! 😄

---

## 📄 Licença

MIT © 2026 **Muriel Gasparini** — veja [LICENSE](LICENSE).
