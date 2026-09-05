<p align="center">
  <img src="https://img.shields.io/badge/Feito%20com-Go-00ADD8?style=plastic&logo=go&logoColor=white" alt="Go"/>
  <img src="https://img.shields.io/badge/vers%C3%A3o-v0.2.0-blue?style=plastic" alt="versão"/>
  <img src="https://img.shields.io/badge/status-est%C3%A1vel-success?style=plastic" alt="status"/>
</p>

<h1 align="center">🚀 Antigravity Account Switcher</h1>

<p align="center"><b>Seu “nunca mais fique sem cota” para o Google Antigravity 2.0.</b></p>

<p align="center">
  <a href="README.md">🇺🇸 English</a> · 🇧🇷 Português (Brasil)
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Plataforma-Windows_%7C_Linux_%7C_macOS-6f42c1?style=plastic&logo=windows&logoColor=white" alt="plataforma"/>
  <img src="https://img.shields.io/github/go-mod/go-version/samucamg/antigravity-account-switcher?style=plastic&logo=go&label=Go" alt="go version"/>
  <img src="https://img.shields.io/github/license/samucamg/antigravity-account-switcher?style=plastic&color=blue" alt="license"/>
  <img src="https://img.shields.io/github/stars/samucamg/antigravity-account-switcher?style=plastic&logo=github" alt="stars"/>
  <img src="https://github.com/samucamg/antigravity-account-switcher/actions/workflows/ci.yml/badge.svg" alt="CI"/>
</p>

<p align="center">
  🧠 Claude &nbsp;·&nbsp; 🤖 GPT &nbsp;·&nbsp; ✨ Gemini — <b>todos sob um pool inteligente de múltiplas contas</b>
</p>

---

## 😩 Você já passou por isso?

Você está no meio de uma sessão incrível de programação no Google Antigravity 2.0… e, do nada:

```text
HTTP 429 RESOURCE_EXHAUSTED
```

Cota esgotada. Fluxo interrompido. Raciocínio perdido. 😤

O **Antigravity Account Switcher** existe para que isso nunca mais aconteça com você.

---

## ✨ O que ele faz por você?

O switcher é um **supervisor local transparente e de latência zero** que fica entre o Antigravity 2.0 e o Google Cloud Code:

| 💜 Recurso | O que acontece |
|:---|:---|
| 🔔 **Alerta aos 80%** | Quando uma conta atinge 80% do uso da cota, você recebe um aviso instantâneo em tempo real (dashboard web + eventos) — dá tempo de agir com calma. |
| 🔁 **Troca proativa aos 85%** | Antes de estourar a cota e travar seu trabalho, o switcher rotaciona suavemente para a próxima conta saudável do seu pool. |
| 🧠 **Fallback Multinível de Modelos** *(v0.2.0)* | Esgotou a cota do modelo premium (ex: Gemini 2.5 Pro)? O proxy reescreve a requisição em tempo real para um modelo secundário (ex: Gemini 2.5 Flash) na **mesma conta** antes de trocá-la! |
| 🔄 **Failover reativo (HTTP 429)** | Se mesmo assim vier um 429, a requisição em voo é repetida em memória com a nova conta — sem janelas de erro no editor. |
| 🚇 **Túneis Cloudflare** *(v0.2.0)* | Tunelamento com 1 clique integrado: crie túneis públicos instantâneos Quick Tunnels (`trycloudflare.com`) ou conecte túneis fixos com token do Cloudflare Zero Trust para acesso remoto seguro. |
| ⚡ **Iniciador Windows em 1 Clique** *(v0.2.0)* | Script `start.bat` e atalho na Área de Trabalho que iniciam o proxy em segundo plano, abrem a dashboard web e já lançam o Antigravity IDE. |
| 📊 **Cotas e Métricas com Fuso Horário** | Monitora o consumo por modelo (Claude, GPT, Gemini) e exibe gráficos de uso mapeados perfeitamente no seu horário local. |
| 🔒 **100% local e privado** | Nenhum token ou prompt sai da sua máquina. Tudo fica protegido no seu banco de dados local SQLite (`accounts.db`). |
| 🪟🐧🍏 **Windows, Linux e macOS** | Um único binário Go de alta performance, sem dependências externas. |

---

## 🛠️ Como funciona (em 10 segundos)

```text
  Google Antigravity 2.0 / agy
          │  (requisições interceptadas pelo proxy local na porta 1831 / 8080)
          ▼
  🧲 Antigravity Account Switcher
          │  utiliza a conta ativa do pool
          │
          ├─ cota ≥ 80%              → 🔔 dispara alerta preventivo
          ├─ cota ≥ 85%              → 🔁 rotação proativa de conta
          ├─ cota do modelo esgotada → 🧠 fallback para modelo secundário (Pro → Flash)
          └─ resposta HTTP 429       → 🔄 troca de conta + repetição em memória
          ▼
  ☁️ Google Cloud Code (Cloud Code PA)
```

---

## 📦 Guia de Instalação

### 🪟 Windows (Instalação Rápida e Recomendada)

Disponibilizamos um instalador automatizado e um inicializador de 1 clique projetados especificamente para usuários do Windows.

#### 1. Pré-requisitos
- **Go 1.24+**: Baixe o instalador MSI para Windows em [go.dev/dl](https://go.dev/dl/) e instale. Confirme abrindo o PowerShell e digitando `go version`.
- **Git para Windows**: Baixe em [git-scm.com](https://git-scm.com/) (ou instale com `winget install Git.Git`).
- **Google Antigravity 2.0 IDE**: Instalado no seu computador (o caminho padrão é `%LOCALAPPDATA%\Programs\Antigravity\Antigravity.exe`).
- *(Opcional - para Túneis Cloudflare)*: Se desejar acesso remoto via túnel, instale o `cloudflared` pelo PowerShell:
  ```powershell
  winget install --id Cloudflare.cloudflared
  ```

---

#### 2. Instalação Automatizada em 1 Passo (`install.ps1`)

Abra o PowerShell, navegue até a pasta onde deseja manter a ferramenta e execute:

```powershell
git clone https://github.com/samucamg/antigravity-account-switcher.git
cd antigravity-account-switcher
powershell -ExecutionPolicy Bypass -File .\install.ps1
```

> [!TIP]
> Se o PowerShell exibir o erro *“a execução de scripts foi desabilitada neste sistema”*, você pode liberar a execução local digitando:
> ```powershell
> Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
> ```
> E depois execute novamente `.\install.ps1`.

**O que o `install.ps1` faz automaticamente para você:**
1. ✅ Verifica e valida o compilador Go.
2. 🔨 Compila o `antigravity-account-switcher.exe` de forma otimizada.
3. ⚙️ Define a porta padrão do serviço para `1831`.
4. 🔍 Localiza automaticamente o executável do Google Antigravity 2.0 IDE.
5. 📄 Cria scripts auxiliares: `start.ps1`, `add-account.ps1`, `switch-account.ps1`.
6. 🖥️ **Cria um Atalho na Área de Trabalho:** Adiciona o atalho **`Antigravity (Multi-Account).lnk`** direto no seu Desktop com o ícone oficial do Antigravity IDE!

---

#### 3. Uso Diário no Windows (1 Clique)

Você tem duas maneiras práticas de começar seu dia de trabalho:

- **Método A (Mais Fácil):** Dê um duplo clique no atalho **`Antigravity (Multi-Account)`** na sua Área de Trabalho!
- **Método B:** Dê um duplo clique no arquivo **`start.bat`** dentro da pasta do projeto.

**O que acontece por trás dos panos:**
- 🛡️ Encerra processos antigos do switcher para evitar conflitos de porta.
- 🚀 Inicia o servidor proxy de forma invisível em segundo plano.
- 🌐 Abre seu navegador padrão diretamente na Dashboard Web em `http://127.0.0.1:1831/`.
- 🔗 Injeta automaticamente as configurações de proxy (`HTTP_PROXY`, `HTTPS_PROXY`, `CLOUD_CODE_URL`) e abre o **Google Antigravity 2.0 IDE**!

---

#### 4. Uso Manual via Terminal no Windows (PowerShell / CMD)

Se você prefere utilizar o terminal diretamente:

```powershell
# Compilar o binário
go build -o antigravity-account-switcher.exe ./cmd/antigravity-account-switcher

# Adicionar conta Google (abre navegador para OAuth2 seguro)
.\antigravity-account-switcher.exe add-account

# Iniciar o Antigravity IDE sob o proxy supervisionado
.\antigravity-account-switcher.exe launch

# Ou rodar apenas o servidor proxy e a dashboard web
.\antigravity-account-switcher.exe serve --port 1831
```

---

#### 5. Resolução de Dúvidas e Problemas no Windows

- **Porta 1831 já está em uso?**
  O `start.bat` encerra automaticamente instâncias anteriores. Se outro aplicativo estiver usando a porta 1831, mude para outra a qualquer momento:
  ```powershell
  .\antigravity-account-switcher.exe config set port 1832
  ```
- **Antigravity IDE instalado em pasta personalizada?**
  Caso o `install.ps1` não tenha localizado a IDE automaticamente, informe o caminho:
  ```powershell
  .\antigravity-account-switcher.exe config set antigravity_bin "D:\Programas\Antigravity\Antigravity.exe"
  ```
- **Aviso do Firewall do Windows:**
  Na primeira inicialização, o Windows pode perguntar sobre permissão de rede para `antigravity-account-switcher.exe`. Selecione **Redes privadas** e clique em **Permitir acesso** para que o Antigravity IDE consiga falar com o proxy local.
- **Usando WSL2 (Windows Subsystem for Linux)?**
  Você pode rodar o switcher inteiramente dentro do WSL2! Siga os passos de Linux abaixo. O switcher detecta e aciona sozinho o executável Windows do Antigravity.

---

### 🐧 Linux / 🍏 macOS

```bash
# 1. Clonar repositório
git clone https://github.com/samucamg/antigravity-account-switcher.git
cd antigravity-account-switcher

# 2. Compilar e instalar binário (instala em ~/.local/bin)
make install

# 3. (Opcional no Linux) Criar ícone no menu de aplicativos
antigravity-account-switcher install-desktop
```

---

## 🏃 Comece em 3 Passos

1️⃣ **Adicione suas contas Google** (quantas quiser — 2, 3, 10…):

```bash
# No Windows
.\antigravity-account-switcher.exe add-account

# No Linux / macOS
antigravity-account-switcher add-account
```
*(Ou clique no botão **“Add Account”** diretamente na Dashboard Web!)*

2️⃣ **Inicie o Antigravity sob supervisão:**
- No **Windows**: Dê duplo clique no atalho **Antigravity (Multi-Account)** na Área de Trabalho ou no `start.bat`.
- No **Linux / macOS**: Execute `antigravity-account-switcher launch`.

3️⃣ **Acompanhe tudo pela Dashboard:**
Abra `http://127.0.0.1:1831/` (ou `http://127.0.0.1:8080/`) no seu navegador para ver contas, velocímetros de cota, gráficos de consumo e logs de troca em tempo real. 👀

> 💡 **Auto-Import:** Já está logado no Google Antigravity? Na primeira execução, o switcher **importa suas credenciais ativas automaticamente** — sem nenhuma configuração manual.

---

## 🌟 Novidades da Versão v0.2.0

### 🧠 Fallback Multinível de Modelos
Durante sessões intensas de código, modelos mais pesados (como `gemini-2.5-pro` ou Claude 3.5 Sonnet) podem esgotar a cota por hora enquanto modelos mais leves na *mesma* conta continuam com 100% livre.
- **Failover inteligente intra-conta:** Ao receber erro de cota esgotada (HTTP 429 / 403 `RESOURCE_EXHAUSTED`), o switcher reescreve a requisição em tempo real para um modelo secundário (ex: `gemini-2.5-flash`) na **mesma conta**, sem precisar queimar as cotas das outras contas do seu pool!
- **Reescrita com alocação zero:** Mecanismo de reescrita em stream de alta performance que ajusta cabeçalhos e corpo JSON em trânsito com impacto nulo na CPU ou memória.
- **Configurável via Web UI ou CLI:**
  ```bash
  antigravity-account-switcher serve --fallback-secondary --model-primary gemini-2.5-pro --model-secondary gemini-2.5-flash
  ```

---

### 🚇 Túneis Cloudflare (Acesso Remoto e Webhooks)
Precisa acessar o painel de cotas do Antigravity Switcher no seu celular, tablet ou em outro computador?
- **Quick Tunnels (1 Clique):** Inicie um túnel HTTPS público e instantâneo via `trycloudflare.com`. Sem necessidade de criar conta no Cloudflare nem cadastrar cartão.
- **Túneis Nomeados do Zero Trust:** Conecte um domínio próprio seguro usando o token de túnel do Cloudflare Zero Trust (`eyJh...`), com políticas de acesso e autenticação corporativa.
- **Integrado à Dashboard Web:** Ative/desative com 1 clique, visualize o status da conexão ao vivo e copie a URL pública com facilidade.
- **Pré-requisito:** Tenha o `cloudflared` instalado (`winget install --id Cloudflare.cloudflared` no Windows ou `brew install cloudflared` no macOS).

---

## 🖥️ Dashboard Web

Acesse a interface interativa em `http://127.0.0.1:1831/`:

- 📋 **Visão Geral do Pool:** Status de cada conta (Ativa, Disponível, Esgotada, Em Recuperação).
- 🧮 **Cotas por Modelo:** Barras de porcentagem em tempo real para Gemini, Claude e GPT.
- 🧠 **Roteamento de Fallback de Modelos:** Ative o fallback intra-conta e escolha modelos primário e secundário com autodescoberta.
- 🚇 **Gerenciamento de Túneis Cloudflare:** Inicie Quick Tunnels ou configure túneis Zero Trust em 1 clique.
- 🔔 **Alertas Visuais:** Avisos destacados quando a cota atinge 80%.
- 📈 **Gráficos Ajustados ao Fuso Horário:** Histórico de consumo de tokens por hora e dia exibidos no seu relógio local.
- ➕ **Adicionar Contas via Web:** Conecte novas contas Google pelo próprio navegador com apenas um clique.

---

## ⌨️ Referência de Comandos da CLI

| Comando | Descrição |
|:---|:---|
| `launch` | **(Recomendado)** Inicia o Antigravity 2.0 sob o proxy supervisionado |
| `serve` | Roda o motor proxy + monitor de cotas + dashboard web |
| `add-account` | Adiciona uma nova conta Google via OAuth2 (loopback + PKCE) |
| `list-accounts` | Lista todas as contas, a conta ativa e o % de cota restante |
| `switch-account [email]` | Alterna manualmente a conta ativa do pool |
| `refresh-quotas` | Força a sincronização imediata de cotas com o Google |
| `wrap -- <cmd>` | Executa qualquer comando injetando as variáveis de ambiente do proxy |
| `status` | Mostra a conta ativa, saúde dos tokens e métricas do supervisor |
| `config` | Visualiza e altera parâmetros de configuração (`get`, `set`, `list`) |
| `install-desktop` | Cria atalho com ícone no menu de aplicativos do Linux |
| `uninstall-desktop` | Remove o atalho do menu de aplicativos do Linux |
| `version` | Exibe a versão (`v0.2.0`), hash do commit e data de compilação |

### Opções (Flags) para `serve` e `launch`
- `--port <porta>`: Porta TCP para escuta (ex: `1831` ou `8080`).
- `--fallback-secondary`: Ativa o fallback para o modelo secundário na mesma conta antes de trocar de conta.
- `--model-primary <id>`: Identificador do modelo primário (padrão: `gemini-2.5-pro` ou autodetectado).
- `--model-secondary <id>`: Identificador do modelo de fallback (padrão: `gemini-2.5-flash`).

---

## ⚙️ Configuração

As configurações são salvas em:
- **Windows:** `%APPDATA%\antigravity-account-switcher\config.json`
- **Linux / macOS:** `~/.config/antigravity-account-switcher/config.json`

```bash
# Ver todas as configurações
antigravity-account-switcher config list

# Ajustar parâmetros
antigravity-account-switcher config set port 1831
antigravity-account-switcher config set quota_interval 60s
antigravity-account-switcher config set antigravity_bin "C:\Caminho\Para\Antigravity.exe"

# Limiares de cota (valores de 0.0 a 1.0)
antigravity-account-switcher config set quota_warning_threshold 0.80   # 🔔 alerta aos 80%
antigravity-account-switcher config set quota_switch_threshold  0.85   # 🔁 troca aos 85%
```

**Variáveis de Ambiente:** `ANTIGRAVITY_BIN`, `ANTIGRAVITY_PORT`, `ANTIGRAVITY_DB_PATH`, `ANTIGRAVITY_CLIENT_ID`, `ANTIGRAVITY_CLIENT_SECRET`.

---

## 🔒 Segurança & Privacidade

- 🛡️ **Padrão IETF RFC OAuth2:** Fluxo loopback seguro com PKCE (RFC 7636) — padrão ouro para autenticação de aplicativos nativos.
- 🏠 **Armazenamento 100% Local:** Credenciais, tokens de atualização e metadados ficam guardados exclusivamente no seu banco SQLite local (`accounts.db`) com permissões restritas.
- 🚫 **Zero Telemetria:** Nenhum prompt, token, código, métrica ou dado pessoal é enviado a terceiros.
- 🔐 **Tunelamento Transparente:** Tráfego que não é do Cloud Code (como reconhecimento de voz em `speech.googleapis.com`) passa por tunelamento TCP `CONNECT` bidirecional sem inspeção.

Consulte [SECURITY.md](SECURITY.md) para obter mais detalhes.

---

## 🤝 Contribuindo

Contribuições são muito bem-vindas! Veja o arquivo [CONTRIBUTING.pt-BR.md](CONTRIBUTING.pt-BR.md) para configurar seu ambiente de desenvolvimento, rodar os testes com verificação de corrida (*race detector*) e abrir seu pull request.

```bash
# Rodar suíte de testes com detecção de concorrência
go test -race ./...
```

---

## ⭐ Gostou do projeto?

Deixe uma ⭐ no repositório do GitHub e compartilhe com outros desenvolvedores que também sofrem com limites de cota! 😄

---

## 📄 Licença

MIT © 2026 **Muriel Gasparini** — veja [LICENSE](LICENSE).
