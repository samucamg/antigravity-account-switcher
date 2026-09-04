# 🚀 Antigravity Account Switcher - Windows Installer & Setup Script

$ErrorActionPreference = "Stop"

Write-Host "`n========================================================" -ForegroundColor Cyan
Write-Host "  🚀 Antigravity Account Switcher - Windows Installer" -ForegroundColor Cyan
Write-Host "========================================================`n" -ForegroundColor Cyan

# 1. Check Go environment
Write-Host "[1/4] Checking Go environment..." -ForegroundColor Yellow
try {
    $goVersion = & go version
    Write-Host "  ✓ Found: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "  ❌ Go toolchain not found in PATH!" -ForegroundColor Red
    Write-Host "  Please install Go 1.24+ from https://go.dev/dl/ and restart PowerShell." -ForegroundColor Yellow
    exit 1
}

# 2. Compile binary
Write-Host "`n[2/4] Compiling antigravity-account-switcher.exe..." -ForegroundColor Yellow
try {
    & go build -o antigravity-account-switcher.exe ./cmd/antigravity-account-switcher
    Write-Host "  ✓ Successfully compiled antigravity-account-switcher.exe" -ForegroundColor Green
} catch {
    Write-Host "  ❌ Compilation failed!" -ForegroundColor Red
    Write-Host $_ -ForegroundColor Red
    exit 1
}

# 3. Configure default port to 1831 (avoids common 8080 conflicts)
Write-Host "`n[3/4] Configuring default port (1831)..." -ForegroundColor Yellow
try {
    & .\antigravity-account-switcher.exe config set port 1831 | Out-Null
    Write-Host "  ✓ Default port set to 1831" -ForegroundColor Green
} catch {
    Write-Host "  ⚠ Notice: Could not set default port automatically." -ForegroundColor Yellow
}

# 4. Create helper scripts (start.ps1 and add-account.ps1)
Write-Host "`n[4/4] Creating 1-click execution scripts..." -ForegroundColor Yellow

$startCode = '$exe = Join-Path $PSScriptRoot "antigravity-account-switcher.exe"; if (Test-Path $exe) { & $exe launch --open } else { Write-Host "Executable not found! Run .\install.ps1 first." -ForegroundColor Red }'
$addCode   = '$exe = Join-Path $PSScriptRoot "antigravity-account-switcher.exe"; if (Test-Path $exe) { & $exe add-account } else { Write-Host "Executable not found! Run .\install.ps1 first." -ForegroundColor Red }'

Set-Content -Path (Join-Path $PSScriptRoot "start.ps1") -Value $startCode -Encoding UTF8
Set-Content -Path (Join-Path $PSScriptRoot "add-account.ps1") -Value $addCode -Encoding UTF8

Write-Host "  ✓ Created .\start.ps1       (Launches Switcher + Antigravity IDE)" -ForegroundColor Green
Write-Host "  ✓ Created .\add-account.ps1 (Add Google Account via Browser)" -ForegroundColor Green

Write-Host "`n========================================================" -ForegroundColor Cyan
Write-Host "  🎉 Setup Complete!" -ForegroundColor Green
Write-Host "========================================================" -ForegroundColor Cyan
Write-Host "`nNext Steps:" -ForegroundColor Yellow
Write-Host "  1. Add your Google accounts:" -ForegroundColor White
Write-Host "     .\add-account.ps1" -ForegroundColor Cyan
Write-Host "`n  2. Launch Antigravity IDE under supervision:" -ForegroundColor White
Write-Host "     .\start.ps1`n" -ForegroundColor Cyan