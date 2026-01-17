@echo off
chcp 65001 >nul
echo.
echo ========================================
echo 物联网卡管理系统 - 快速启动
echo ========================================
echo.

REM 检查Go
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo 未安装 Go，请先安装 Go 1.21+
    pause
    exit /b 1
)

REM 检查Node.js
where node >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo 未安装 Node.js，请先安装 Node.js 18+
    pause
    exit /b 1
)

echo 依赖检查完成
echo.

echo ========================================
echo 数据库配置说明
echo ========================================
echo.
echo 请确保PostgreSQL数据库已经配置：
echo.
echo 方式1：使用Docker（推荐）
echo   docker run -d --name iot_postgres -e POSTGRES_DB=iot_card_db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres123 -p 5432:5432 postgres:15-alpine
echo.
echo 方式2：使用本地PostgreSQL
echo   确保PostgreSQL服务正在运行
echo.
echo 然后执行数据库迁移：
echo   docker exec -i iot_postgres psql -U postgres -d iot_card_db ^< backend/migrations/001_create_admin_users.sql
echo   docker exec -i iot_postgres psql -U postgres -d iot_card_db ^< backend/migrations/002_create_sim_cards.sql
echo   docker exec -i iot_postgres psql -U postgres -d iot_card_db ^< backend/migrations/003_create_recharge_records.sql
echo   docker exec -i iot_postgres psql -U postgres -d iot_card_db ^< backend/migrations/004_create_system_config.sql
echo   docker exec -i iot_postgres psql -U postgres -d iot_card_db ^< backend/migrations/insert_test_data.sql
echo.
echo ========================================
echo.
pause

echo.
echo [1/3] 启动后端服务...
cd backend
start "IoT后端" cmd /k go run cmd/server/main.go
cd ..
timeout /t 5 /nobreak >nul
echo 后端启动完成 (http://localhost:8080)
echo.

echo [2/3] 启动H5前端...
cd frontend\h5
if not exist node_modules (
    echo 正在安装H5依赖，请稍候...
    call npm install
)
start "IoT-H5" cmd /k npm run dev
cd ..\..
timeout /t 3 /nobreak >nul
echo H5启动完成 (http://localhost:3000)
echo.

echo [3/3] 启动管理端...
cd frontend\admin
if not exist node_modules (
    echo 正在安装管理端依赖，请稍候...
    call npm install
)
start "IoT-Admin" cmd /k npm run dev
cd ..\..
timeout /t 3 /nobreak >nul
echo 管理端启动完成 (http://localhost:3001)
echo.

echo ========================================
echo 所有服务启动完成！
echo ========================================
echo.
echo H5端: http://localhost:3000
echo   测试卡号: 89860123456789012345
echo.
echo 管理端: http://localhost:3001
echo   账号: admin
echo   密码: admin123
echo.
echo 后端API: http://localhost:8080/api/v1
echo.
echo 每个服务都在独立窗口运行
echo 关闭对应窗口即可停止服务
echo ========================================
pause
