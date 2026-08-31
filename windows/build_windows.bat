@echo off
echo ===================================================
echo     SingUI Windows Native Build Script
echo ===================================================

cd frontend
call npm install
call npm run build
cd ..

set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0
go build -ldflags="-s -w" -o singbox-ui.exe main.go

echo ===================================================
echo   Build Completed: singbox-ui.exe
echo ===================================================
pause
