# Antigravity Account Switcher - Start Script
# Starts the proxy server in background, opens the web dashboard,
# then launches Antigravity IDE with the proxy environment injected.

$ErrorActionPreference = "Stop"
$exe  = Join-Path $PSScriptRoot "antigravity-account-switcher.exe"
$port = 1831

if (-not (Test-Path $exe)) {
    Write-Host "ERROR: exe not found. Run .\install.ps1 first." -ForegroundColor Red
    exit 1
}

# Kill leftover server from previous session
Get-Process -Name "antigravity-account-switcher" -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
Start-Sleep -Milliseconds 400

# Start proxy + dashboard in hidden background window
Write-Host "Starting proxy server on port $port..." -ForegroundColor Cyan
Start-Process -FilePath $exe -ArgumentList "serve","--port",$port -WindowStyle Hidden

# Wait up to 10s for server to respond
Write-Host "Waiting for server..." -ForegroundColor Cyan
$ready = $false
for ($i = 0; $i -lt 20; $i++) {
    Start-Sleep -Milliseconds 500
    try {
        $r = Invoke-WebRequest -Uri "http://127.0.0.1:$port/api/status" -UseBasicParsing -TimeoutSec 1 -ErrorAction SilentlyContinue
        if ($r.StatusCode -eq 200) { $ready = $true; break }
    } catch {}
}
if (-not $ready) {
    Write-Host "ERROR: Server did not start. Check for port conflicts." -ForegroundColor Red
    exit 1
}

Write-Host "Server running at http://127.0.0.1:$port/" -ForegroundColor Green
Start-Process "http://127.0.0.1:$($port)/"
Write-Host "Dashboard opened in browser." -ForegroundColor Green

# Launch IDE with proxy environment
$ideBin = "C:\Users\samuc\AppData\Local\Programs\Antigravity IDE\Antigravity IDE.exe"

if ($ideBin -eq '' -or -not (Test-Path $ideBin)) {
    Write-Host "IDE not found. Edit start.ps1 and set the ideBin path." -ForegroundColor Yellow
    Write-Host "Server still running at http://127.0.0.1:$port/" -ForegroundColor Cyan
    exit 0
}

$env:HTTP_PROXY     = "http://127.0.0.1:$port"
$env:HTTPS_PROXY    = "http://127.0.0.1:$port"
$env:http_proxy     = "http://127.0.0.1:$port"
$env:https_proxy    = "http://127.0.0.1:$port"
$env:NO_PROXY       = "localhost,127.0.0.1,::1"
$env:no_proxy       = "localhost,127.0.0.1,::1"
$env:CLOUD_CODE_URL = "http://127.0.0.1:$port"

Write-Host "Launching Antigravity IDE through proxy..." -ForegroundColor Cyan
& $ideBin
Write-Host "IDE closed. Server still running at http://127.0.0.1:$port/" -ForegroundColor Cyan
