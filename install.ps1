# Antigravity Account Switcher - Windows Installer & Setup Script

$ErrorActionPreference = "Stop"

Write-Host ""
Write-Host "========================================================" -ForegroundColor Cyan
Write-Host "  Antigravity Account Switcher - Windows Installer" -ForegroundColor Cyan
Write-Host "========================================================" -ForegroundColor Cyan
Write-Host ""

# 0. Stop any background switcher instance
Get-Process -Name "antigravity-account-switcher" -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue

# 1. Check Go environment
Write-Host "[1/5] Checking Go environment..." -ForegroundColor Yellow
$goVersion = & go version 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "  ERROR: Go not found in PATH. Install from https://go.dev/dl/" -ForegroundColor Red; exit 1
}
Write-Host "  OK: $goVersion" -ForegroundColor Green

# 2. Compile
Write-Host ""
Write-Host "[2/5] Compiling antigravity-account-switcher.exe..." -ForegroundColor Yellow
$buildOut = & go build -o antigravity-account-switcher.exe ./cmd/antigravity-account-switcher 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "  ERROR: $buildOut" -ForegroundColor Red; exit 1
}
Write-Host "  OK: compiled." -ForegroundColor Green

# 3. Set default port
Write-Host ""
Write-Host "[3/5] Configuring default port (1831)..." -ForegroundColor Yellow
& .\antigravity-account-switcher.exe config set port 1831 2>&1 | Out-Null
Write-Host "  OK: port = 1831" -ForegroundColor Green

# 4. Detect IDE
Write-Host ""
Write-Host "[4/5] Detecting Antigravity IDE..." -ForegroundColor Yellow
$idePaths = @(
    "$env:LOCALAPPDATA\Programs\Antigravity\Antigravity.exe",
    "$env:LOCALAPPDATA\Programs\Antigravityntigravity.exe",
    "$env:LOCALAPPDATA\Programs\Antigravity IDE\Antigravity IDE.exe",
    "C:\Program Files\Antigravity\Antigravity.exe",
    "C:\Program Files\Antigravity IDE\Antigravity IDE.exe"
)
$detectedIDE = ""
foreach ($p in $idePaths) { if (Test-Path $p) { $detectedIDE = $p; break } }
if ($detectedIDE -ne "") {
    Write-Host "  OK: $detectedIDE" -ForegroundColor Green
} else {
    Write-Host "  WARN: IDE not found in standard paths. Edit start.ps1 after setup." -ForegroundColor Yellow
}

# 5. Generate helper scripts
Write-Host ""
Write-Host "[5/5] Creating scripts..." -ForegroundColor Yellow

$enc = [System.Text.Encoding]::UTF8

# start.ps1
$startContent = @"
# Antigravity Account Switcher - Start Script
`$exe  = Join-Path `$PSScriptRoot 'antigravity-account-switcher.exe'
`$port = 1831
if (-not (Test-Path `$exe)) { Write-Host 'ERROR: Run .\install.ps1 first.' -ForegroundColor Red; exit 1 }
Get-Process -Name 'antigravity-account-switcher' -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
Start-Sleep -Milliseconds 400
Write-Host 'Starting proxy server on port ' -NoNewline; Write-Host `$port -ForegroundColor Cyan
Start-Process -FilePath `$exe -ArgumentList 'serve','--port',`$port -WindowStyle Hidden
Write-Host 'Waiting for server...' -ForegroundColor Cyan
`$ready = `$false
for (`$i = 0; `$i -lt 20; `$i++) {
    Start-Sleep -Milliseconds 500
    try { `$r = Invoke-WebRequest -Uri "http://127.0.0.1:`$port/api/status" -UseBasicParsing -TimeoutSec 1 -EA SilentlyContinue; if (`$r.StatusCode -eq 200) { `$ready = `$true; break } } catch {}
}
if (-not `$ready) { Write-Host 'ERROR: Server did not start. Check port conflicts.' -ForegroundColor Red; exit 1 }
Write-Host "Server running - Dashboard: http://127.0.0.1:`$port/" -ForegroundColor Green
Start-Process "http://127.0.0.1:`$(`$port)/"
`$ideBin = '$detectedIDE'
if (`$ideBin -ne '' -and (Test-Path `$ideBin)) {
    `$env:HTTP_PROXY = "http://127.0.0.1:`$port"; `$env:HTTPS_PROXY = "http://127.0.0.1:`$port"
    `$env:http_proxy = "http://127.0.0.1:`$port"; `$env:https_proxy = "http://127.0.0.1:`$port"
    `$env:NO_PROXY = 'localhost,127.0.0.1,::1'; `$env:no_proxy = 'localhost,127.0.0.1,::1'
    `$env:CLOUD_CODE_URL = "http://127.0.0.1:`$port"
    Write-Host 'Launching Antigravity IDE through proxy...' -ForegroundColor Cyan
    & `$ideBin
    Write-Host 'IDE closed. Server still running.' -ForegroundColor Cyan
} else {
    Write-Host 'IDE not found. Edit start.ps1 and set `$ideBin.' -ForegroundColor Yellow
}
"@
[System.IO.File]::WriteAllText((Join-Path $PSScriptRoot 'start.ps1'), $startContent, $enc)
Write-Host "  OK: start.ps1" -ForegroundColor Green

# add-account.ps1
$addContent = '& (Join-Path $PSScriptRoot "antigravity-account-switcher.exe") add-account'
[System.IO.File]::WriteAllText((Join-Path $PSScriptRoot 'add-account.ps1'), $addContent, $enc)
Write-Host "  OK: add-account.ps1" -ForegroundColor Green

# switch-account.ps1
$switchContent = @"
# Usage: .\switch-account.ps1              -> lists accounts
#        .\switch-account.ps1 email@x.com  -> sets active account
`$exe = Join-Path `$PSScriptRoot 'antigravity-account-switcher.exe'
if (`$args.Count -eq 0) { & `$exe list-accounts } else { & `$exe switch-account @args }
"@
[System.IO.File]::WriteAllText((Join-Path $PSScriptRoot 'switch-account.ps1'), $switchContent, $enc)
Write-Host "  OK: switch-account.ps1" -ForegroundColor Green

# Create Desktop Shortcut for 1-Click Launch
try {
    $desktopDir = [Environment]::GetFolderPath("Desktop")
    $shortcutPath = Join-Path $desktopDir "Antigravity (Multi-Account).lnk"
    $wsh = New-Object -ComObject WScript.Shell
    $shortcut = $wsh.CreateShortcut($shortcutPath)
    $shortcut.TargetPath = Join-Path $PSScriptRoot "start.bat"
    $shortcut.WorkingDirectory = $PSScriptRoot
    $shortcut.Description = "Iniciar Google Antigravity 2.0 com Proxy Multi-Account"
    if ($detectedIDE -ne '' -and (Test-Path $detectedIDE)) {
        $shortcut.IconLocation = "$detectedIDE,0"
    }
    $shortcut.Save()
    Write-Host "  OK: Atalho criado na Area de Trabalho ('Antigravity (Multi-Account).lnk')" -ForegroundColor Green
} catch {
    Write-Host "  AVISO: Nao foi possivel criar o atalho na Area de Trabalho: $_" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "========================================================" -ForegroundColor Cyan
Write-Host "  Setup Complete!" -ForegroundColor Green
Write-Host "========================================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "  start.bat                            -> 1-click startup (or double-click Desktop shortcut)" -ForegroundColor Cyan
Write-Host "  .\start.ps1                          -> start server + dashboard + launch IDE" -ForegroundColor Cyan
Write-Host "  .\add-account.ps1                    -> add Google account" -ForegroundColor Cyan
Write-Host "  .\switch-account.ps1                 -> list accounts" -ForegroundColor Cyan
Write-Host "  .\switch-account.ps1 email@gmail.com -> switch active account" -ForegroundColor Cyan
Write-Host ""
