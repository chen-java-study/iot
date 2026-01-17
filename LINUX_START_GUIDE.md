# Linux 物联网卡管理系统启动指南

## 🚀 快速开始

### 方法一：一键启动脚本（推荐）

```bash
# 克隆项目
git clone <repository-url>
cd iot-card-management

# 给脚本执行权限
chmod +x start-dev.sh

# 一键启动所有服务
./start-dev.sh
```

### 方法二：手动分步启动

## 📋 前置要求

- **Go 1.21+**
- **Node.js 18+**
- **Docker & Docker Compose** (推荐)
- **PostgreSQL 15+** (可选，如果不使用Docker)

### Ubuntu/Debian 安装依赖

```bash
# 更新包管理器
sudo apt update

# 安装Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# 安装Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# 安装Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

## 🗄️ 第一步：启动数据库

### 方式A：使用Docker（推荐）

```bash
# 启动PostgreSQL容器
docker run -d \
  --name iot_postgres \
  -e POSTGRES_DB=iot_card_db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres123 \
  -p 5432:5432 \
  postgres:15-alpine

# 等待数据库启动
sleep 10

# 执行数据库迁移
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend/migrations/001_create_admin_users.sql
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend/migrations/002_create_sim_cards.sql
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend/migrations/003_create_recharge_records.sql
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend/migrations/004_create_system_config.sql
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend/migrations/insert_test_data.sql
```

### 方式B：使用本地PostgreSQL

```bash
# 安装PostgreSQL
sudo apt install postgresql postgresql-contrib

# 启动服务
sudo systemctl start postgresql
sudo systemctl enable postgresql

# 创建数据库
sudo -u postgres createdb iot_card_db

# 执行迁移脚本
sudo -u postgres psql -d iot_card_db -f backend/migrations/001_create_admin_users.sql
sudo -u postgres psql -d iot_card_db -f backend/migrations/002_create_sim_cards.sql
sudo -u postgres psql -d iot_card_db -f backend/migrations/003_create_recharge_records.sql
sudo -u postgres psql -d iot_card_db -f backend/migrations/004_create_system_config.sql
sudo -u postgres psql -d iot_card_db -f backend/migrations/insert_test_data.sql
```

## 🔧 第二步：启动后端服务

```bash
# 进入后端目录
cd backend

# 下载Go依赖
go mod download

# 启动后端服务（后台运行）
go run cmd/server/main.go &

# 验证后端启动
sleep 3
curl http://localhost:8080/api/v1/card/query?keyword=test
```

**成功标志：**
```
[GIN-debug] Listening and serving HTTP on :8080
```

## 📱 第三步：启动H5前端

```bash
# 进入H5前端目录
cd frontend/h5

# 安装依赖
npm install

# 启动开发服务器（后台运行）
npm run dev &

# 或者指定端口
npm run dev -- --port 3000 &
```

**成功标志：**
```
VITE v5.x.x  ready in xxx ms
➜  Local:   http://localhost:3000/
```

## 💼 第四步：启动管理端（可选）

```bash
# 进入管理端目录
cd frontend/admin

# 安装依赖
npm install

# 启动开发服务器
npm run dev &
```

## 🌐 访问地址

启动完成后，访问以下地址：

- **H5端**: http://localhost:3000
- **管理端**: http://localhost:3001
- **后端API**: http://localhost:8080/api/v1

## 🧪 测试验证

### 测试H5端
1. 打开浏览器访问 http://localhost:3000
2. 输入测试卡号：`89860123456789012345`
3. 点击"查询"，应该能看到卡片信息

### 测试管理端
1. 打开浏览器访问 http://localhost:3001
2. 账号：`admin`
3. 密码：`admin123`
4. 登录后应该能看到统计数据和卡片管理

## 🛑 停止服务

```bash
# 查看进程
ps aux | grep -E "(go|vite|node)"

# 杀死进程（根据PID）
kill -9 <PID>

# 或者使用pkill
pkill -f "go run cmd/server/main.go"
pkill -f "npm run dev"
```

## 🔧 故障排除

### 问题1：后端启动失败

```bash
# 检查数据库连接
docker ps | grep iot_postgres

# 查看后端日志
cd backend && go run cmd/server/main.go

# 测试数据库连接
psql -h localhost -U postgres -d iot_card_db
```

### 问题2：前端启动失败

```bash
# 使用国内npm镜像
npm config set registry https://registry.npmmirror.com

# 清除缓存重新安装
rm -rf node_modules package-lock.json
npm install
```

### 问题3：端口被占用

```bash
# 查看端口占用
lsof -i :8080
lsof -i :3000
lsof -i :3001

# 修改端口
# 后端：编辑 backend/configs/config.yaml
# 前端：编辑 frontend/h5/vite.config.js 或 frontend/admin/vite.config.js
```

## 🚀 生产环境部署

### 使用Docker Compose

```bash
# 构建前端
cd frontend/h5 && npm install && npm run build
cd ../admin && npm install && npm run build

# 启动所有服务
cd ../..
docker-compose up -d

# 访问地址
# H5: http://localhost/h5
# 管理端: http://localhost/admin
# API: http://localhost/api
```

### 查看日志

```bash
# Docker Compose日志
docker-compose logs -f

# 单独查看服务日志
docker-compose logs -f backend
docker-compose logs -f nginx
```

## 📊 测试数据

系统预置了5张测试卡片：

1. `89860123456789012345` - 正常状态（中国移动）
2. `89860123456789012346` - 即将到期（中国联通）
3. `89860123456789012347` - 已过期（中国电信）
4. `89860123456789012348` - 正常状态（中国移动）
5. `89860123456789012349` - 正常状态（中国联通）

## 🔧 常用命令

```bash
# 检查服务状态
docker ps
docker-compose ps

# 重启服务
docker-compose restart

# 停止所有服务
docker-compose down

# 清理数据
docker-compose down -v
docker system prune -f
```
