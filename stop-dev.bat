@echo off
chcp 65001 >nul
echo.
echo ========================================
echo 🛑 停止所有服务
echo ========================================
echo.

echo 🔍 查找并停止Go进程...
taskkill /F /FI "WINDOWTITLE eq 物联网卡-后端*" >nul 2>nul
if %ERRORLEVEL% EQU 0 (
    echo ✅ 后端服务已停止
) else (
    echo ⚠️  未找到运行中的后端服务
)

echo.
echo 🔍 查找并停止H5进程...
taskkill /F /FI "WINDOWTITLE eq 物联网卡-H5*" >nul 2>nul
if %ERRORLEVEL% EQU 0 (
    echo ✅ H5服务已停止
) else (
    echo ⚠️  未找到运行中的H5服务
)

echo.
echo 🔍 查找并停止管理端进程...
taskkill /F /FI "WINDOWTITLE eq 物联网卡-管理端*" >nul 2>nul
if %ERRORLEVEL% EQU 0 (
    echo ✅ 管理端服务已停止
) else (
    echo ⚠️  未找到运行中的管理端服务
)

echo.
echo ========================================
echo ✅ 停止完成
echo ========================================
pause
