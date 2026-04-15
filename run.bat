@echo off
REM PortManager Desktop GUI
REM 自动以管理员权限运行

echo.
echo ======================================
echo   PortManager - 端口管理工具
echo ======================================
echo.

REM 检查是否以管理员身份运行
net session >nul 2>&1
if %errorlevel% neq 0 (
    echo 正在请求管理员权限...
    powershell -Command "Start-Process cmd -ArgumentList '/c \"%~f0\"' -Verb RunAs"
    exit /b
)

REM 运行程序
echo 正在启动 PortManager...
"%~dp0PortManager.exe"
