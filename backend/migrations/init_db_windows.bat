@echo off
REM 数据库初始化脚本 - Windows
REM 自动创建数据库、用户和所有表

echo ========================================
echo 数据库初始化脚本
echo ========================================

REM 检查PostgreSQL是否安装
where psql >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo [错误] 未找到PostgreSQL命令行工具(psql)
    echo 请确保PostgreSQL已安装且添加到系统PATH中
    pause
    exit /b 1
)

echo.
echo [步骤1] 创建数据库和用户...
psql -U postgres -f "%~dp0\00_init_database.sql"
if %ERRORLEVEL% neq 0 (
    echo [错误] 数据库初始化失败
    pause
    exit /b 1
)

echo.
echo [步骤2] 创建管理员用户表...
psql -U iot_user -h localhost -d iot_card_db -f "%~dp0\001_create_admin_users.sql"
if %ERRORLEVEL% neq 0 (
    echo [错误] 创建admin_users表失败
    pause
    exit /b 1
)

echo.
echo [步骤3] 创建SIM卡表...
psql -U iot_user -h localhost -d iot_card_db -f "%~dp0\002_create_sim_cards.sql"
if %ERRORLEVEL% neq 0 (
    echo [错误] 创建sim_cards表失败
    pause
    exit /b 1
)

echo.
echo [步骤4] 创建充值记录表...
psql -U iot_user -h localhost -d iot_card_db -f "%~dp0\003_create_recharge_records.sql"
if %ERRORLEVEL% neq 0 (
    echo [错误] 创建recharge_records表失败
    pause
    exit /b 1
)

echo.
echo [步骤5] 创建系统配置表...
psql -U iot_user -h localhost -d iot_card_db -f "%~dp0\004_create_system_config.sql"
if %ERRORLEVEL% neq 0 (
    echo [错误] 创建system_config表失败
    pause
    exit /b 1
)

echo.
echo [步骤6] 插入测试数据...
psql -U iot_user -h localhost -d iot_card_db -f "%~dp0\insert_test_data.sql"
if %ERRORLEVEL% neq 0 (
    echo [警告] 插入测试数据失败，但表已成功创建
)

echo.
echo ========================================
echo 数据库初始化完成！
echo ========================================
echo.
echo 数据库信息:
echo   - 数据库名: iot_card_db
echo   - 用户名: iot_user
echo   - 密码: adfhkIxcvYIK2189
echo   - 主机: localhost
echo   - 端口: 5432
echo.
echo 已创建的表:
echo   - admin_users (管理员用户表)
echo   - sim_cards (SIM卡信息表)
echo   - recharge_records (充值记录表)
echo   - system_config (系统配置表)
echo.
pause
