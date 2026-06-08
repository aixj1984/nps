@echo off
setlocal
cd /d "%~dp0"
set GOPROXY=https://proxy.golang.org,direct

echo [1/4] build frontend...
cd frontend
call npm install
if errorlevel 1 exit /b 1
call npm run build
if errorlevel 1 exit /b 1
cd ..

echo [2/4] go mod tidy...
go mod tidy
if errorlevel 1 exit /b 1

echo [3/4] wails generate module...
where wails >nul 2>&1
if %errorlevel%==0 (
  wails generate module
) else (
  echo WARN: wails CLI not found, skipping binding generation
)

echo [4/4] go build...
go build -tags desktop,production -ldflags "-w -s -H windowsgui" -o npc-desk.exe .
if errorlevel 1 exit /b 1

echo Done: npc-desk.exe
endlocal
