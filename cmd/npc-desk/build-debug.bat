@echo off
setlocal
cd /d "%~dp0"
echo Manual DEV build — see https://wails.io/docs/guides/manual-builds/
if not exist frontend\dist\index.html (
  echo Run build.bat first.
  exit /b 1
)
go run ./cmd/genwindows
if errorlevel 1 exit /b 1
go build -tags dev -gcflags "all=-N -l" -ldflags "-H windowsgui" -o npc-deskreen-debug.exe .
if errorlevel 1 exit /b 1
if exist npc-deskreen-res.syso del /q npc-deskreen-res.syso
echo Done: npc-deskreen-debug.exe
endlocal
