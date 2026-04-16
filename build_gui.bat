@echo off
setlocal
cd /d "%~dp0"

REM 生成资源文件（如果rsrc工具可用）
if exist "%GOPATH%\bin\rsrc.exe" (
    echo 正在嵌入图标和manifest...
    "%GOPATH%\bin\rsrc.exe" -manifest PortManager.exe.manifest -ico PortManager.ico -o rsrc_windows_amd64.syso
    if errorlevel 1 (
        echo 警告：资源文件生成失败，继续使用现有资源
    )
)

echo 正在编译应用...
go build -ldflags "-H=windowsgui" -o PortManager.exe .
if errorlevel 1 (
    echo.
    echo 构建失败。
    pause
    exit /b 1
)
echo.
echo 构建完成：PortManager.exe (windowsgui 含图标)
pause
