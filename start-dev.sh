#!/bin/bash

# 物联网卡管理系统 - 快速测试脚本

echo "🚀 开始启动物联网卡管理系统..."
echo ""

# 检查依赖
check_dependencies() {
    echo "📋 检查依赖..."

    if ! command -v go &> /dev/null; then
        echo "❌ 未安装 Go，请先安装 Go 1.21+"
        exit 1
    fi

    if ! command -v node &> /dev/null; then
        echo "❌ 未安装 Node.js，请先安装 Node.js 18+"
        exit 1
    fi

    if ! command -v docker &> /dev/null; then
        echo "⚠️  未安装 Docker，将尝试连接本地数据库"
        USE_DOCKER=false
    else
        USE_DOCKER=true
    fi

    echo "✅ 依赖检查完成"
    echo ""
}

# 启动数据库
start_database() {
    echo "🗄️  启动数据库..."

    if [ "$USE_DOCKER" = true ]; then
        # 检查容器是否已存在
        if docker ps -a | grep -q iot_postgres; then
            echo "📦 PostgreSQL容器已存在，启动中..."
            docker start iot_postgres
        else
            echo "📦 创建并启动PostgreSQL容器..."
            docker run -d \
              --name iot_postgres \
              -e POSTGRES_DB=iot_card_db \
              -e POSTGRES_USER=postgres \
              -e POSTGRES_PASSWORD=postgres123 \
              -p 5432:5432 \
              postgres:15-alpine

            echo "⏳ 等待数据库启动..."
            sleep 10

            echo "📝 执行数据库迁移..."
            for sql_file in backend/migrations/*.sql; do
                echo "   执行 $sql_file"
                docker exec -i iot_postgres psql -U postgres -d iot_card_db < "$sql_file"
            done
        fi
    else
        echo "⚠️  请确保PostgreSQL已安装并运行在 localhost:5432"
        echo "   数据库名: iot_card_db"
        echo "   用户名: postgres"
    fi

    echo "✅ 数据库准备完成"
    echo ""
}

# 启动后端
start_backend() {
    echo "🔧 启动后端服务..."

    cd backend

    # 安装依赖
    if [ ! -f "go.sum" ]; then
        echo "📦 下载Go依赖..."
        go mod download
    fi

    # 启动服务
    echo "🚀 启动Go服务器..."
    go run cmd/server/main.go &
    BACKEND_PID=$!

    echo "   后端PID: $BACKEND_PID"
    echo "   访问地址: http://localhost:8080"

    cd ..

    # 等待后端启动
    echo "⏳ 等待后端启动..."
    sleep 5

    # 测试后端是否正常
    if curl -s http://localhost:8080/api/v1/card/query?keyword=test > /dev/null; then
        echo "✅ 后端启动成功"
    else
        echo "⚠️  后端可能启动失败，请检查日志"
    fi

    echo ""
}

# 启动H5前端
start_h5() {
    echo "📱 启动H5前端..."

    cd frontend/h5

    # 安装依赖
    if [ ! -d "node_modules" ]; then
        echo "📦 安装npm依赖..."
        npm install
    fi

    # 启动开发服务器
    echo "🚀 启动H5开发服务器..."
    npm run dev > ../../logs/h5.log 2>&1 &
    H5_PID=$!

    echo "   H5 PID: $H5_PID"
    echo "   访问地址: http://localhost:3000"

    cd ../..
    echo "✅ H5前端启动完成"
    echo ""
}

# 启动管理端
start_admin() {
    echo "💼 启动管理端..."

    cd frontend/admin

    # 安装依赖
    if [ ! -d "node_modules" ]; then
        echo "📦 安装npm依赖..."
        npm install
    fi

    # 启动开发服务器
    echo "🚀 启动管理端开发服务器..."
    npm run dev > ../../logs/admin.log 2>&1 &
    ADMIN_PID=$!

    echo "   管理端PID: $ADMIN_PID"
    echo "   访问地址: http://localhost:3001"

    cd ../..
    echo "✅ 管理端启动完成"
    echo ""
}

# 显示测试信息
show_info() {
    echo "════════════════════════════════════════"
    echo "✅ 所有服务启动完成！"
    echo "════════════════════════════════════════"
    echo ""
    echo "📱 H5端测试："
    echo "   访问: http://localhost:3000"
    echo "   测试卡号: 89860123456789012345"
    echo ""
    echo "💼 管理端测试："
    echo "   访问: http://localhost:3001"
    echo "   账号: admin"
    echo "   密码: admin123"
    echo ""
    echo "🔧 后端API："
    echo "   地址: http://localhost:8080/api/v1"
    echo ""
    echo "📊 查看日志："
    echo "   tail -f logs/h5.log"
    echo "   tail -f logs/admin.log"
    echo ""
    echo "🛑 停止服务："
    echo "   ./scripts/stop.sh"
    echo "   或按 Ctrl+C"
    echo ""
    echo "════════════════════════════════════════"
}

# 创建日志目录
mkdir -p logs

# 执行启动流程
check_dependencies
start_database
start_backend
start_h5
start_admin
show_info

# 保存PIDs
echo $BACKEND_PID > logs/backend.pid
echo $H5_PID > logs/h5.pid
echo $ADMIN_PID > logs/admin.pid

# 等待用户中断
echo "💡 提示: 按 Ctrl+C 停止所有服务"
wait
