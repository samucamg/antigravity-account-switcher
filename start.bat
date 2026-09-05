@echo off
setlocal enabledelayedexpansion
title Antigravity Launcher

cd /d "%~dp0"

echo [1/4] Verificando executavel do Switcher...
if not exist "antigravity-account-switcher.exe" (
    echo ERRO: antigravity-account-switcher.exe nao encontrado nesta pasta!
    echo Execute install.ps1 primeiro.
    pause
    exit /b 1
)

echo [2/4] Detectando Google Antigravity 2.0 IDE...
set "IDE_BIN="

if exist "%LOCALAPPDATA%\Programs\Antigravity\Antigravity.exe" (
    set "IDE_BIN=%LOCALAPPDATA%\Programs\Antigravity\Antigravity.exe"
) else if exist "%LOCALAPPDATA%\Programs\Antigravity IDE\Antigravity IDE.exe" (
    set "IDE_BIN=%LOCALAPPDATA%\Programs\Antigravity IDE\Antigravity IDE.exe"
) else if exist "%ProgramFiles%\Antigravity\Antigravity.exe" (
    set "IDE_BIN=%ProgramFiles%\Antigravity\Antigravity.exe"
) else if exist "%ProgramFiles%\Antigravity IDE\Antigravity IDE.exe" (
    set "IDE_BIN=%ProgramFiles%\Antigravity IDE\Antigravity IDE.exe"
) else if exist "%ProgramFiles(x86)%\Antigravity\Antigravity.exe" (
    set "IDE_BIN=%ProgramFiles(x86)%\Antigravity\Antigravity.exe"
)

if "%IDE_BIN%"=="" (
    echo AVISO: Antigravity.exe nao encontrado nos caminhos padroes.
    echo O servidor proxy sera iniciado, mas a IDE deve ser aberta manualmente.
) else (
    echo Encontrado: "%IDE_BIN%"
)

echo [3/4] Encerrando processos antigos e iniciando proxy em segundo plano (porta 1831)...
taskkill /F /IM antigravity-account-switcher.exe >nul 2>&1
timeout /t 1 /nobreak >nul

start "" /B antigravity-account-switcher.exe serve --port 1831

echo [4/4] Abrindo Dashboard e iniciando Antigravity IDE...
timeout /t 2 /nobreak >nul
start "" "http://127.0.0.1:1831/"

if not "%IDE_BIN%"=="" (
    set "HTTP_PROXY=http://127.0.0.1:1831"
    set "HTTPS_PROXY=http://127.0.0.1:1831"
    set "http_proxy=http://127.0.0.1:1831"
    set "https_proxy=http://127.0.0.1:1831"
    set "NO_PROXY=localhost,127.0.0.1,::1"
    set "no_proxy=localhost,127.0.0.1,::1"
    set "CLOUD_CODE_URL=http://127.0.0.1:1831"
    start "" "%IDE_BIN%"
)

exit /b 0
