@echo off
set GOOS= linux
set GOARCH=amd64
go build -o srvimglinux release.go

REM Confirmar si la compilación fue exitosa
if %errorlevel% neq 0 (
    echo Error en la compilación
    exit /b %errorlevel%
)

echo Compilación completada con éxito. El ejecutable es nombrearchivo-linux
pause
