# 物联网卡管理系统

一个完整的物联网SIM卡管理系统，包含手机端H5查询充值和Web管理后台。

## 目录

1. [项目概述](#项目概述)
2. [技术架构](#技术架构)
3. [功能特性](#功能特性)
4. [快速开始](#快速开始)
5. [本地开发](#本地开发)
6. [数据库配置](#数据库配置)
7. [API文档](#api文档)
8. [阿里云部署](#阿里云部署)
9. [项目结构](#项目结构)
10. [安全建议](#安全建议)
11. [故障排查](#故障排查)
12. [缺陷修复记录](#缺陷修复记录)

---

## 项目概述

### 系统组成

- **手机端H5**: 用户查询SIM卡到期时间、在线充值（微信支付）
- **管理端Web**: 商家管理卡片、查看充值记录、系统配置

### 业务流程

1. 用户在H5端输入卡号/设备号查询卡片信息
2. 查看到期时间，点击充值跳转支付页面
3. 完成微信支付后，卡片到期时间自动延长1年
4. 商家在管理端录入卡片、查看充值记录、对账

---

## 技术架构

### 技术栈

| 层级 | 技术 |
|------|------|
| **后端** | Go 1.21 + Gin框架 + GORM |
| **数据库** | PostgreSQL 15 |
| **认证** | JWT |
| **支付** | 微信支付APIv3 |
| **H5前端** | Vue 3 + Vite + Vant 4 |
| **管理前端** | Vue 3 + Vite + Element Plus |
| **部署** | Docker + Nginx |

### 架构图

```
┌─────────────┐     ┌─────────────┐
│   H5前端    │     │  管理端前端  │
│  (Vant 4)   │     │(Element Plus)│
└──────┬──────┘     └──────┬──────┘
       │                   │
       └─────────┬─────────┘
                 │
         ┌───────▼───────┐
         │    Nginx      │
         │  (反向代理)    │
         └───────┬───────┘
                 │
         ┌───────▼───────┐
         │   Go Backend  │
         │   (Gin框架)   │
         └───────┬───────┘
                 │
         ┌───────▼───────┐
         │  PostgreSQL   │
         └───────────────┘
```

---

## 功能特性

### 手机端H5

- 📱 查询SIM卡信息（输入卡号/设备号）
- 💳 在线充值（微信支付）
- 📊 查看到期时间和剩余天数
- ✅ 充值后自动延长1年

### 管理端Web

- 🏠 数据仪表盘（卡片总数、充值总额统计）
- 📋 卡片管理（增删改查）
- 💰 充值记录查询
- ⚙️ 系统配置（充值价格、微信支付配置）
- 🔐 用户登录认证（JWT）

---

## 快速开始

### 前置要求

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+ 或 Docker
- 微信支付商户号（生产环境）

### 方式一：一键启动（推荐）

**Windows:**
```cmd
双击运行：start-dev.bat
```

**Linux/Mac:**
```bash
chmod +x start-dev.sh
./start-dev.sh
```

### 方式二：Docker快速启动

```bash
# 1. 启动PostgreSQL



# 3. 执行数据库迁移
for sql in backend/migrations/00*.sql backend/migrations/insert_test_data.sql; do
  docker exec -i iot_postgres psql -U postgres -d iot_card_db < "$sql"
done

# 4. 启动后端
cd backend && go run cmd/server/main.go &

# 5. 启动H5前端
cd frontend/h5 && npm install && npm run dev &

# 6. 启动管理端
cd frontend/admin && npm install && npm run dev &
```

### 方式三：生产环境部署（Systemd服务）

```bash
# 1. 编译后端程序
cd /home/workspace/iot/backend
GOMAXPROCS=1 CGO_ENABLED=0 go build -o iot-server ./cmd/server/

# 2. 创建日志目录
mkdir -p /home/workspace/iot/logs

# 3. 安装systemd服务
sudo cp /home/workspace/iot/backend/iot-backend.service /etc/systemd/system/
sudo systemctl daemon-reload

# 4. 启动并设置开机自启
sudo systemctl start iot-backend
sudo systemctl enable iot-backend
```

生产环境服务管理:
```bash
sudo systemctl status iot-backend    # 查看状态
sudo systemctl restart iot-backend   # 重启服务
sudo journalctl -u iot-backend -f    # 实时日志
```

### 访问地址

| 服务 | 地址 |
|------|------|
| H5端 | http://localhost:3000 |
| 管理端 | http://localhost:3001 |
| 后端API | http://localhost:8080/api/v1 |

### 测试账号

- **管理员**: admin / admin123
- **测试卡号**: 89860123456789012345

---



**生产环境（Systemd服务）：**
```bash
# 1. 编译二进制程序
cd /home/workspace/iot/backend
GOMAXPROCS=1 CGO_ENABLED=0 go build -o iot-server ./cmd/server/

# 2. 创建日志目录
mkdir -p /home/workspace/iot/logs
export DB_PASSWORD="adfhkIxcvYIK2189"
nohup ./iot-server > /home/workspace/iot/logs/iot-server.log 2>&1 &
# 3. 安装并启动服务
sudo cp /home/workspace/iot/backend/iot-backend.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl start iot-backend
sudo systemctl enable iot-backend
```

生产环境服务管理：
```bash
sudo systemctl status iot-backend     # 查看状态
sudo systemctl restart iot-backend    # 重启服务
sudo journalctl -u iot-backend -f     # 实时日志
```

#### 4. 启动H5前端

```bash
cd frontend/h5
npm install
npm run dev
```

成功标志：
```
VITE v5.x.x ready in xxx ms
➜ Local: http://localhost:3000/
```

#### 5. 启动管理端

```bash
cd frontend/admin
npm install
npm run dev

cd /home/iot/iot-master/frontend/admin
nohup npm run dev > /home/workspace/iot/logs/vue.log 2>&1 &
echo "前端已启动 (PID: $!)"
```

成功标志：
```
VITE v5.x.x ready in xxx ms
➜ Local: http://localhost:3001/
```

---

## 数据库配置

### 表结构说明

| 表名 | 说明 |
|------|------|
| admin_users | 管理员用户表 |
| sim_cards | SIM卡信息表 |
| recharge_records | 充值记录表 |
| system_config | 系统配置表 |
| v_statistics | 统计视图 |

### 数据库凭证

| 项目 | 开发环境 | 生产环境 |
|------|----------|----------|
| 主机 | localhost | 使用环境变量 DB_HOST |
| 端口 | 5432 | 5432 |
| 数据库 | iot_card_db | iot_card_db |
| 用户 | postgres | 使用环境变量 DB_USER |
| 密码 | postgres123 | 使用环境变量 DB_PASSWORD |

### 配置文件

编辑 `backend/configs/config.yaml`:

```yaml
server:
  port: 8080
  mode: release

database:
  host: 127.0.0.1
  port: 5432
  user: iot_user
  password: ""          # 生产环境使用环境变量 DB_PASSWORD
  dbname: iot_card_db

jwt:
  secret_key: ""        # 生产环境使用环境变量 JWT_SECRET
  expire_hours: 24

wechat:
  app_id: ""            # 生产环境使用环境变量 WECHAT_APP_ID
  mch_id: ""            # 生产环境使用环境变量 WECHAT_MCH_ID
  api_v3_key: ""        # 生产环境使用环境变量 WECHAT_API_V3_KEY
  serial_no: ""         # 生产环境使用环境变量 WECHAT_SERIAL_NO
  private_key_path: "/path/to/apiclient_key.pem"
  notify_url: "https://yourdomain.com/api/v1/payment/notify"
```

---

## API文档

### H5端API（无需认证）

**查询卡片**
```
GET /api/v1/card/query?keyword=卡号或设备号

响应:
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "card_no": "89860123456789012345",
    "device_no": "866123456789012",
    "expire_date": "2027-01-01",
    "status": 1,
    "operator": "中国移动"
  }
}
```

**创建充值订单**
```
POST /api/v1/payment/create
Content-Type: application/json

{
  "card_no": "89860123456789012345",
  "openid": "user_wechat_openid"
}
```

**查询订单状态**
```
GET /api/v1/payment/status?trade_no=订单号
```

### 管理端API（需JWT认证）

**登录**
```
POST /api/v1/admin/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}

响应:
{
  "code": 200,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": { "id": 1, "username": "admin" }
  }
}
```

**获取统计数据**
```
GET /api/v1/admin/statistics
Headers: Authorization: Bearer <token>
```

**卡片列表**
```
GET /api/v1/admin/cards?page=1&page_size=20&status=1&keyword=
Headers: Authorization: Bearer <token>
```

**创建卡片**
```
POST /api/v1/admin/cards
Headers: Authorization: Bearer <token>
Content-Type: application/json

{
  "card_no": "89860123456789012345",
  "device_no": "866123456789012",
  "start_date": "2026-01-01T00:00:00Z",
  "expire_date": "2027-01-01T00:00:00Z",
  "operator": "中国移动",
  "package_type": "年卡"
}
```

**更新卡片**
```
PUT /api/v1/admin/cards/:id
Headers: Authorization: Bearer <token>
```

**删除卡片**
```
DELETE /api/v1/admin/cards/:id
Headers: Authorization: Bearer <token>
```

**充值记录列表**
```
GET /api/v1/admin/recharges?page=1&page_size=20&status=1&keyword=&start_date=&end_date=
Headers: Authorization: Bearer <token>
```

**获取系统配置**
```
GET /api/v1/admin/config
Headers: Authorization: Bearer <token>
```

**更新系统配置**
```
POST /api/v1/admin/config
Headers: Authorization: Bearer <token>
Content-Type: application/json

{
  "recharge_price": "100.00",
  "wechat_app_id": "your_app_id"
}
```

---

## 阿里云部署

### 资源准备

| 资源 | 推荐配置 |
|------|----------|
| ECS云服务器 | ecs.c6.large (2核4G), Ubuntu 22.04 |
| RDS PostgreSQL | rds.pg.s2.large (2核4G), 50GB SSD |
| 域名 | 需完成ICP备案 |
| SSL证书 | 阿里云免费DV证书 |

### 部署步骤

#### 1. 服务器环境配置

```bash
# 安装Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# 安装Nginx
sudo apt install nginx -y
```

#### 2. 数据库配置

```sql
-- 连接RDS后执行
CREATE DATABASE iot_card_db WITH ENCODING 'UTF8';
CREATE USER iot_user WITH PASSWORD '你的强密码';
GRANT ALL PRIVILEGES ON DATABASE iot_card_db TO iot_user;
```

#### 3. 编译部署后端

```bash
cd /opt/iot-card-system/backend
go mod tidy
GOMAXPROCS=1 CGO_ENABLED=0 go build -o iot-server ./cmd/server/
```

#### 4. 构建部署前端

```bash
# H5前端
cd /opt/iot-card-system/frontend/h5
npm install && npm run build
cp -r dist/* /var/www/iot-h5/

# 管理端
cd /opt/iot-card-system/frontend/admin
npm install && npm run build
cp -r dist/* /var/www/iot-admin/
```

#### 5. Nginx配置

创建 `/etc/nginx/conf.d/iot-card.conf`:

```nginx
server {
    listen 443 ssl http2;
    server_name admin.yourdomain.com;

    ssl_certificate /etc/nginx/ssl/admin.yourdomain.com.pem;
    ssl_certificate_key /etc/nginx/ssl/admin.yourdomain.com.key;

    root /var/www/iot-admin;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}

server {
    listen 443 ssl http2;
    server_name h5.yourdomain.com;

    ssl_certificate /etc/nginx/ssl/h5.yourdomain.com.pem;
    ssl_certificate_key /etc/nginx/ssl/h5.yourdomain.com.key;

    root /var/www/iot-h5;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

#### 6. 环境变量配置

创建 `/opt/iot-card-system/backend/.env`:

```bash
DB_HOST=rm-xxxxx.pg.rds.aliyuncs.com
DB_USER=iot_user
DB_PASSWORD=你的数据库密码
DB_NAME=iot_card_db
JWT_SECRET=生成一个32位以上的随机字符串
WECHAT_APP_ID=你的AppID
WECHAT_MCH_ID=你的商户号
WECHAT_API_V3_KEY=你的APIv3密钥
WECHAT_SERIAL_NO=商户证书序列号
WECHAT_PRIVATE_KEY_PATH=/opt/iot-card-system/certs/apiclient_key.pem
WECHAT_NOTIFY_URL=https://h5.yourdomain.com/api/v1/payment/notify
```

#### 7. Systemd服务配置

创建 `/etc/systemd/system/iot-backend.service`:

```ini
[Unit]
Description=IoT Backend Service
After=network.target postgresql.service
Wants=postgresql.service

[Service]
Type=simple
User=root
WorkingDirectory=/home/workspace/iot/backend
ExecStart=/home/workspace/iot/backend/iot-server
Restart=always
RestartSec=5

# 日志配置
StandardOutput=append:/home/workspace/iot/logs/backend.log
StandardError=append:/home/workspace/iot/logs/backend_error.log

[Install]
WantedBy=multi-user.target
```

安装并启动服务:
```bash
# 1. 编译后端程序
cd /home/workspace/iot/backend
GOMAXPROCS=1 CGO_ENABLED=0 go build -o iot-server ./cmd/server/

# 2. 创建日志目录
mkdir -p /home/workspace/iot/logs

# 3. 安装systemd服务
sudo cp /home/workspace/iot/backend/iot-backend.service /etc/systemd/system/
sudo systemctl daemon-reload

# 4. 启动并设置开机自启
sudo systemctl start iot-backend
sudo systemctl enable iot-backend
```

服务管理命令:
```bash
sudo systemctl status iot-backend    # 查看状态
sudo systemctl restart iot-backend   # 重启服务
sudo systemctl stop iot-backend      # 停止服务
sudo journalctl -u iot-backend -f    # 实时日志
```

### 部署检查清单

- [ ] ECS服务器已购买并配置安全组
- [ ] RDS数据库已创建并配置白名单
- [ ] 数据库表结构已初始化
- [ ] 域名已备案并解析到ECS
- [ ] SSL证书已申请并配置
- [ ] 后端代码已编译并部署
- [ ] 前端代码已构建并部署
- [ ] Nginx配置已完成
- [ ] 环境变量已配置
- [ ] Systemd服务已配置并启动 (iot-backend)
- [ ] 微信支付证书已上传

---

## 项目结构

```
iot-card-management/
├── backend/                    # Go后端
│   ├── cmd/server/main.go      # 程序入口
│   ├── internal/
│   │   ├── config/             # 配置管理
│   │   ├── handler/            # HTTP处理器
│   │   ├── middleware/         # 中间件(CORS, JWT)
│   │   ├── model/              # 数据模型
│   │   ├── repository/         # 数据访问层
│   │   ├── router/             # 路由配置
│   │   ├── service/            # 业务逻辑层
│   │   └── utils/              # 工具函数
│   ├── migrations/             # 数据库迁移脚本
│   ├── pkg/database/           # 数据库连接
│   ├── configs/                # 配置文件
│   └── go.mod
├── frontend/
│   ├── h5/                     # 手机端H5 (Vue3 + Vant4)
│   │   ├── src/
│   │   │   ├── api/            # API接口
│   │   │   ├── views/          # 页面组件
│   │   │   ├── router/         # 路由配置
│   │   │   └── utils/          # 工具函数
│   │   └── vite.config.js
│   └── admin/                  # 管理端 (Vue3 + Element Plus)
│       ├── src/
│       │   ├── api/
│       │   ├── views/
│       │   ├── router/
│       │   └── utils/
│       └── vite.config.js
├── nginx/                      # Nginx配置
├── docker-compose.yml          # Docker编排
└── README.md
```

---

## 安全建议

- ✅ 使用环境变量管理敏感配置（数据库密码、JWT密钥、微信支付密钥）
- ✅ 定期更新JWT密钥
- ✅ 启用HTTPS传输
- ✅ 验证所有微信支付回调签名
- ✅ 使用bcrypt加密密码
- ✅ 实施SQL注入防护（GORM参数化查询）
- ✅ 实施XSS防护（前端输入过滤）
- ✅ 配置CORS白名单（生产环境）

---

## 故障排查

### 后端启动失败

```bash
# ========== 开发环境 ==========
# 查看详细日志
cd backend && go run cmd/server/main.go

# ========== 生产环境(Systemd服务) ==========
# 查看服务状态
sudo systemctl status iot-backend

# 查看实时日志
sudo journalctl -u iot-backend -f

# 查看应用日志文件
tail -f /home/workspace/iot/logs/backend.log
tail -f /home/workspace/iot/logs/backend_error.log

# 测试数据库连接
psql -h localhost -U postgres -d iot_card_db
```

### 前端启动失败

```bash
# 使用国内npm镜像
npm config set registry https://registry.npmmirror.com

# 清除缓存重新安装
rm -rf node_modules package-lock.json
npm install
```

### 端口被占用

```bash
# Windows
netstat -ano | findstr :8080

# Linux/Mac
lsof -i :8080
lsof -i :3000
lsof -i :3001
```

### CORS跨域错误

检查后端 `internal/middleware/cors.go` 中的配置，确保允许前端域名。

### Systemd服务异常

```bash
# 检查服务状态和错误信息
sudo systemctl status iot-backend

# 查看详细日志
sudo journalctl -u iot-backend -e  # 从末尾开始显示

# 重新加载配置
sudo systemctl daemon-reload
sudo systemctl restart iot-backend

# 检查配置文件语法
sudo systemd-analyze verify /etc/systemd/system/iot-backend.service
```

### 微信支付回调失败

1. 确认域名已备案且HTTPS正常
2. 检查Nginx代理配置
3. 验证微信支付证书配置正确
4. 检查回调URL是否可访问

---

## 缺陷修复记录

### 已修复的问题

| 问题 | 文件 | 修复内容 |
|------|------|----------|
| CORS配置冲突 | `middleware/cors.go` | `AllowOrigins: ["*"]` 与 `AllowCredentials: true` 不能同时使用，改用 `AllowOriginFunc` |
| AdminLogin错误调用 | `service/service.go` | 更新最后登录时间时错误调用了 `UpdateCard`，改为 `UpdateAdminUser` |
| 总金额计算SQL错误 | `repository/repository.go` | `query.Statement.SQL.String()` 方式错误，改为重新构建查询条件 |
| 配置文件敏感信息 | `configs/config.yaml` | 移除硬编码密码，支持环境变量覆盖 |
| 支付回调安全验证 | `utils/wechat.go` | 新增微信支付签名验证工具 |

---

## 测试数据

系统预置了5张测试卡片：

| 卡号 | 状态 | 运营商 |
|------|------|--------|
| 89860123456789012345 | 正常 | 中国移动 |
| 89860123456789012346 | 即将到期 | 中国联通 |
| 89860123456789012347 | 已过期 | 中国电信 |
| 89860123456789012348 | 正常 | 中国移动 |
| 89860123456789012349 | 正常 | 中国联通 |

---

## 许可证

MIT License

---

## 更新日志

### v1.0.2
- 📝 添加生产环境 systemd 服务配置文档
- 📝 添加 systemd 服务异常排查指南
- 📝 完善快速开始和手动分步启动章节

### v1.0.1
- 🐛 修复CORS配置冲突问题
- 🐛 修复AdminLogin更新最后登录时间的错误调用
- 🐛 修复充值记录总金额计算SQL错误
- 🔒 移除配置文件中的硬编码敏感信息
- 🔒 添加环境变量支持
- 🔒 添加微信支付签名验证工具
- 📝 合并整理项目文档

### v1.0.0
- ✨ 初始版本发布
- ✨ 完整的卡片管理功能
- ✨ 微信支付集成
- ✨ 管理后台
- ✨ Docker部署支持

cd /opt/iot-card-system/backend
nohup ./iot-server > ../logs/backend.log 2>&1 &  

cd /opt/iot-card-system/frontend/admin
npm run build
nohup npm run dev > ../../logs/admin.log 2>&1 &

cd /opt/iot-card-system/frontend/h5
npm run build
nohup npm run dev > ../../logs/h5.log 2>&1 &

管理端：http://iot4you.top/admin
H5端：http://iot4you.top/h5

sudo nginx -s reload 


第一步：创建 .env 文件（存放敏感配置）
# 进入项目目录
cd /opt/iot-backend

# 创建 .env 文件
vim .env
# 进入项目目录cd /opt/iot-backend# 创建 .env 文件vim .env
写入以下内容（改成你自己的）：
DB_PASSWORD=your_db_password
WECHAT_APP_ID=wx1234567890abcdef
WECHAT_MCH_ID=1234567890
WECHAT_API_V3_KEY=your_api_v3_key_xxx
WECHAT_SERIAL_NO=ABC123456789
WECHAT_PRIVATE_KEY_PATH=/etc/ssl/wechat/apiclient_key.pem
JWT_SECRET=your_jwt_secret_xxx
DB_PASSWORD=your_db_passwordWECHAT_APP_ID=wx1234567890abcdefWECHAT_MCH_ID=1234567890WECHAT_API_V3_KEY=your_api_v3_key_xxxWECHAT_SERIAL_NO=ABC123456789WECHAT_PRIVATE_KEY_PATH=/etc/ssl/wechat/apiclient_key.pemJWT_SECRET=your_jwt_secret_xxx
第二步：创建 start.sh 启动脚本
vim start.sh
vim start.sh
写入：
#!/bin/bash
# 加载 .env 文件中的环境变量
export $(cat .env | xargs)

# 启动程序
./main
#!/bin/bash# 加载 .env 文件中的环境变量export $(cat .env | xargs)# 启动程序./main
第三步：设置权限
# 只有你能读写 .env（其他人看不到）
chmod 600 .env

# start.sh 可以执行
chmod +x start.sh
# 只有你能读写 .env（其他人看不到）chmod 600 .env# start.sh 可以执行chmod +x start.sh
第四步：启动服务（后台运行）
# 启动（输出日志到 nohup.out）
nohup ./start.sh > nohup.out 2>&1 &

# 说明：
# nohup        不挂断运行（退出终端后继续运行）
# ./start.sh   运行启动脚本
# > nohup.out  标准输出写入 nohup.out 文件
# 2>&1        错误输出也写入同一文件
# &           后台运行
# 启动（输出日志到 nohup.out）nohup ./start.sh > nohup.out 2>&1 &# 说明：# nohup        不挂断运行（退出终端后继续运行）# ./start.sh   运行启动脚本# > nohup.out  标准输出写入 nohup.out 文件# 2>&1        错误输出也写入同一文件# &           后台运行
第五步：验证是否启动成功
# 方式1：看日志
tail -f nohup.out

# 方式2：看进程
ps aux | grep main
# 应该能看到类似：
# www-data 12345  0.0  0.5 123456  5432 ?        S    14:30   0:00 ./main

# 方式3：看端口（如果监听 8080）
netstat -tlnp | grep 8080
# 或
ss -tlnp | grep 8080
# 方式1：看日志tail -f nohup.out# 方式2：看进程ps aux | grep main# 应该能看到类似：# www-data 12345  0.0  0.5 123456  5432 ?        S    14:30   0:00 ./main# 方式3：看端口（如果监听 8080）netstat -tlnp | grep 8080# 或ss -tlnp | grep 8080
第六步：停止服务
# 找到进程 ID
ps aux | grep ./main

# 杀掉进程
kill 12345   # 替换成实际的 PID

# 或者强制杀掉
kill -9 12345
# 找到进程 IDps aux | grep ./main# 杀掉进程kill 12345   # 替换成实际的 PID# 或者强制杀掉kill -9 12345
完整流程演示
# 1. 创建配置
cd /opt/iot-backend
cat > .env << 'EOF'
DB_PASSWORD=mypassword123
WECHAT_APP_ID=wx1234567890
WECHAT_MCH_ID=9876543210
WECHAT_API_V3_KEY=abc123xyz
WECHAT_SERIAL_NO=XYZ123456
WECHAT_PRIVATE_KEY_PATH=/etc/ssl/wechat/apiclient_key.pem
EOF

# 2. 创建启动脚本
cat > start.sh << 'EOF'
#!/bin/bash
export $(cat .env | xargs)
./main
EOF

# 3. 设置权限
chmod 600 .env
chmod +x start.sh

# 4. 启动
nohup ./start.sh > nohup.out 2>&1 &

# 5. 检查
ps aux | grep main
tail -f nohup.out
# 1. 创建配置cd /opt/iot-backendcat > .env << 'EOF'DB_PASSWORD=mypassword123WECHAT_APP_ID=wx1234567890WECHAT_MCH_ID=9876543210WECHAT_API_V3_KEY=abc123xyzWECHAT_SERIAL_NO=XYZ123456WECHAT_PRIVATE_KEY_PATH=/etc/ssl/wechat/apiclient_key.pemEOF# 2. 创建启动脚本cat > start.sh << 'EOF'#!/bin/bashexport $(cat .env | xargs)./mainEOF# 3. 设置权限chmod 600 .envchmod +x start.sh# 4. 启动nohup ./start.sh > nohup.out 2>&1 &# 5. 检查ps aux | grep maintail -f nohup.out
文件结构
/opt/iot-backend/
├── main              # 程序文件
├── start.sh          # 启动脚本
├── .env              # 敏感配置（600 权限）
├── configs/         # 配置文件
│   └── config.yaml
└── nohup.out         # 日志文件（启动后自动生成）
/opt/iot-backend/├── main              # 程序文件├── start.sh          # 启动脚本├── .env              # 敏感配置（600 权限）├── configs/         # 配置文件│   └── config.yaml└── nohup.out         # 日志文件（启动后自动生成）


修改配置步骤
# 1. 编辑 .env 文件
vim /opt/iot-backend/.env

# 2. 修改内容，比如改密码：
DB_PASSWORD=new_password_here
WECHAT_APP_ID=new_app_id
# ... 其他配置

# 3. 保存退出（:wq）

# 4. 重启服务
#    先找到进程
ps aux | grep ./main
#    杀掉
kill 12345

#    重新启动
nohup ./start.sh > nohup.out 2>&1 &
# 1. 编辑 .env 文件vim /opt/iot-backend/.env# 2. 修改内容，比如改密码：DB_PASSWORD=new_password_hereWECHAT_APP_ID=new_app_id# ... 其他配置# 3. 保存退出（:wq）# 4. 重启服务#    先找到进程ps aux | grep ./main#    杀掉kill 12345#    重新启动nohup ./start.sh > nohup.out 2>&1 &
一条命令完成（不进入 vim）
# 用 sed 直接修改（以 DB_PASSWORD 为例）
sed -i 's/DB_PASSWORD=.*/DB_PASSWORD=new_password/' /opt/iot-backend/.env

# 重启
pkill -f ./main && nohup ./start.sh > nohup.out 2>&1 &
# 用 sed 直接修改（以 DB_PASSWORD 为例）sed -i 's/DB_PASSWORD=.*/DB_PASSWORD=new_password/' /opt/iot-backend/.env# 重启pkill -f ./main && nohup ./start.sh > nohup.out 2>&1 &
