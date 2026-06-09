@echo off
setlocal
cd /d "%~dp0"
set GOPROXY=https://proxy.golang.org,direct

echo [1/6] build viewer webui...
cd webui
call npm install
if errorlevel 1 exit /b 1
call npm run build
if errorlevel 1 exit /b 1
cd ..

echo [2/6] build host frontend...
cd frontend
call npm install
if errorlevel 1 exit /b 1
call npm run build
if errorlevel 1 exit /b 1
cd ..

echo [3/6] go mod tidy...
go mod tidy
if errorlevel 1 exit /b 1

if not exist build\appicon.png (
  echo [3b] create default build\appicon.png...
  powershell -NoProfile -Command "Add-Type -AssemblyName System.Drawing; $b=New-Object System.Drawing.Bitmap 256,256; $g=[System.Drawing.Graphics]::FromImage($b); $g.Clear([System.Drawing.Color]::FromArgb(255,16,107,163)); $g.Dispose(); New-Item -Force -ItemType Directory build | Out-Null; $b.Save('build\appicon.png',[System.Drawing.Imaging.ImageFormat]::Png); $b.Dispose()"
  if errorlevel 1 exit /b 1
)

echo [4/6] wails generate module (bindings)...
where wails >nul 2>&1
if %errorlevel%==0 (
  wails generate module
) else (
  echo WARN: wails CLI not found, skipping binding generation
)

echo [5/6] Windows assets (.ico + .syso) - manual build step...
go run ./cmd/genwindows
if errorlevel 1 exit /b 1

echo [6/6] go build -tags desktop,production ...
go build -tags desktop,production -ldflags "-w -s -H windowsgui" -o npc-deskreen.exe .
if errorlevel 1 exit /b 1

if exist npc-deskreen-res.syso del /q npc-deskreen-res.syso

echo Done: npc-deskreen.exe
echo See https://wails.io/docs/guides/manual-builds/
endlocal
