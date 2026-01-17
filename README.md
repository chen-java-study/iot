# 物联网卡管理系统

一个完整的物联网SIM卡管理系统，包含手机端H5查询充值和Web管理后台。

## 项目架构

### 技术栈

**后端:**
- Go 1.21 + Gin框架
- PostgreSQL 15
- JWT认证
- 微信支付APIv3

**前端:**
- Vue 3 + Vite
- H5端: Vant 4
- 管理端: Element Plus
- Axios

**部署:**
- Docker + Docker Compose
- Nginx反向代理

## 功能特性

### 手机端H5 (微信公众号)
- 📱 查询SIM卡信息 (输入卡号/设备号)
- 💳 在线充值 (微信支付)
- 📊 查看到期时间和剩余天数
- ✅ 充值后自动延长1年

### 管理端Web
- 🏠 数据仪表盘 (卡片总数、充值总额统计)
- 📋 卡片管理 (增删改查)
- 💰 充值记录查询
- ⚙️ 系统配置 (充值价格、微信支付配置)
- 🔐 用户登录认证 (JWT)

## 快速开始

### 前置要求

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Docker & Docker Compose (可选)

### 本地开发

#### 1. 克隆项目

\`\`\`bash
git clone <repository-url>
cd iot-card-management
\`\`\`

#### 2. 初始化数据库

\`\`\`bash
# 创建数据库
createdb iot_card_db

# 执行迁移脚本
psql -d iot_card_db -f backend/migrations/001_create_admin_users.sql
psql -d iot_card_db -f backend/migrations/002_create_sim_cards.sql
psql -d iot_card_db -f backend/migrations/003_create_recharge_records.sql
psql -d iot_card_db -f backend/migrations/004_create_system_config.sql
\`\`\`

#### 3. 配置后端

编辑 \`backend/configs/config.yaml\`:

\`\`\`yaml
server:
  port: 8080
  mode: debug

database:
  host: localhost
  port: 5432
  user: postgres
  password: your_password
  dbname: iot_card_db

jwt:
  secret_key: your_jwt_secret_key_change_in_production
  expire_hours: 24

wechat:
  app_id: "your_wechat_appid"
  mch_id: "your_merchant_id"
  api_v3_key: "your_api_v3_key"
  serial_no: "your_certificate_serial_number"
  private_key_path: "path/to/apiclient_key.pem"
  notify_url: "https://yourdomain.com/api/v1/payment/notify"
\`\`\`

#### 4. 启动后端

\`\`\`bash
cd backend
go mod download
go run cmd/server/main.go
\`\`\`

#### 5. 启动H5前端

\`\`\`bash
cd frontend/h5
npm install
npm run dev
# 访问: http://localhost:3000
\`\`\`

#### 6. 启动管理端

\`\`\`bash
cd frontend/admin
npm install
npm run dev
# 访问: http://localhost:3001
\`\`\`

**默认管理员账号:** admin / admin123

### Docker部署

#### 1. 构建并启动所有服务

\`\`\`bash
# 前端项目需先构建
cd frontend/h5
npm install && npm run build

cd ../admin
npm install && npm run build

# 启动Docker服务
cd ../..
docker-compose up -d
\`\`\`

#### 2. 访问系统

- H5端: http://localhost/h5
- 管理端: http://localhost/admin
- API: http://localhost/api

#### 3. 查看日志

\`\`\`bash
docker-compose logs -f backend
\`\`\`

#### 4. 停止服务

\`\`\`bash
docker-compose down
\`\`\`

## API文档

### H5端API (无需认证)

**查询卡片**
\`\`\`
GET /api/v1/card/query?keyword=卡号或设备号
\`\`\`

**创建充值订单**
\`\`\`
POST /api/v1/payment/create
Content-Type: application/json

{
  "card_no": "89860123456789012345",
  "openid": "user_wechat_openid"
}
\`\`\`

**查询订单状态**
\`\`\`
GET /api/v1/payment/status?trade_no=订单号
\`\`\`

### 管理端API (需JWT认证)

**登录**
\`\`\`
POST /api/v1/admin/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
\`\`\`

**获取统计数据**
\`\`\`
GET /api/v1/admin/statistics
Headers: Authorization: Bearer <token>
\`\`\`

**卡片列表**
\`\`\`
GET /api/v1/admin/cards?page=1&page_size=20&status=1&keyword=
Headers: Authorization: Bearer <token>
\`\`\`

**创建卡片**
\`\`\`
POST /api/v1/admin/cards
Headers: Authorization: Bearer <token>
Content-Type: application/json

{
  "card_no": "89860123456789012345",
  "device_no": "866123456789012",
  "start_date": "2026-01-01",
  "expire_date": "2027-01-01",
  "operator": "中国移动",
  "package_type": "年卡"
}
\`\`\`

## 数据库结构

### admin_users (管理员表)
- id, username, password_hash, real_name
- status (1:启用 0:禁用)
- last_login_at, created_at, updated_at

### sim_cards (SIM卡表)
- id, card_no (卡号), device_no (设备号)
- start_date, expire_date
- status (1:正常 2:即将到期 3:已过期)
- operator (运营商), package_type (套餐类型)
- total_recharge_count, total_recharge_amount
- 自动状态更新触发器

### recharge_records (充值记录表)
- id, card_id, card_no, device_no
- recharge_amount, recharge_years
- old_expire_date, new_expire_date
- trade_no (订单号), transaction_id (微信交易号)
- payment_status (0:待支付 1:已支付 2:已退款 3:失败)
- paid_at, openid, ip_address

### system_config (系统配置表)
- id, config_key, config_value
- config_type, description

### v_statistics (统计视图)
- total_cards, active_cards, expiring_cards, expired_cards
- total_recharge_amount, total_recharge_count
- today_amount, month_amount

## 微信支付对接

### 1. 前期准备

- 申请微信公众号(服务号)
- 开通微信支付，获取商户号
- 下载API证书 (apiclient_key.pem)
- 配置支付授权目录
- 设置支付回调地址

### 2. 获取用户OpenID

用户访问H5时，系统会引导进行OAuth授权获取OpenID，用于发起支付。

### 3. 支付流程

1. 前端调用创建订单API
2. 后端调用微信统一下单API
3. 后端返回支付参数
4. 前端调用WeixinJSBridge发起支付
5. 用户完成支付
6. 微信异步回调通知后端
7. 后端验证签名，更新订单，延长卡片到期时间

### 4. 注意事项

- 回调地址必须使用HTTPS
- 需要验证微信回调签名
- 使用事务确保数据一致性
- 防止重复处理支付回调

## 项目结构

\`\`\`
iot-card-management/
├── backend/                    # Go后端
│   ├── cmd/server/main.go      # 程序入口
│   ├── internal/
│   │   ├── config/             # 配置管理
│   │   ├── handler/            # HTTP处理器
│   │   ├── middleware/         # 中间件
│   │   ├── model/              # 数据模型
│   │   ├── repository/         # 数据访问层
│   │   ├── router/             # 路由配置
│   │   ├── service/            # 业务逻辑层
│   │   └── utils/              # 工具函数
│   ├── migrations/             # 数据库迁移
│   ├── pkg/database/           # 数据库连接
│   ├── configs/                # 配置文件
│   ├── go.mod
│   └── Dockerfile
├── frontend/
│   ├── h5/                     # 手机端H5
│   │   ├── src/
│   │   │   ├── api/            # API接口
│   │   │   ├── views/          # 页面组件
│   │   │   ├── router/         # 路由配置
│   │   │   └── utils/          # 工具函数
│   │   ├── package.json
│   │   └── vite.config.js
│   └── admin/                  # 管理端Web
│       ├── src/
│       │   ├── api/
│       │   ├── views/          # 页面组件
│       │   ├── router/
│       │   └── utils/
│       ├── package.json
│       └── vite.config.js
├── nginx/
│   └── nginx.conf              # Nginx配置
├── docker-compose.yml          # Docker编排
└── README.md
\`\`\`

## 开发说明

### 添加新的API接口

1. 在 \`internal/handler/handler.go\` 添加处理函数
2. 在 \`internal/service/service.go\` 添加业务逻辑
3. 在 \`internal/repository/repository.go\` 添加数据访问方法
4. 在 \`internal/router/router.go\` 注册路由

### 添加前端页面

**H5端:**
1. 在 \`frontend/h5/src/views/\` 创建Vue组件
2. 在 \`frontend/h5/src/router/index.js\` 添加路由
3. 在 \`frontend/h5/src/api/\` 添加API调用

**管理端:**
1. 在 \`frontend/admin/src/views/\` 创建Vue组件
2. 在 \`frontend/admin/src/router/index.js\` 添加路由
3. 在 \`frontend/admin/src/api/\` 添加API调用

## 安全建议

- ✅ 使用环境变量管理敏感配置
- ✅ 定期更新JWT密钥
- ✅ 启用HTTPS传输
- ✅ 验证所有微信支付回调签名
- ✅ 使用bcrypt加密密码
- ✅ 实施SQL注入防护(GORM参数化查询)
- ✅ 实施XSS防护(前端输入过滤)

## 性能优化建议

1. **数据库优化**
   - 已添加关键索引
   - 使用数据库连接池
   - 定期清理过期数据

2. **缓存策略**
   - 可添加Redis缓存系统配置
   - 缓存统计数据

3. **前端优化**
   - 路由懒加载
   - 资源压缩和CDN加速

## 故障排查

### 后端启动失败

\`\`\`bash
# 检查数据库连接
psql -h localhost -U postgres -d iot_card_db

# 查看详细日志
cd backend && go run cmd/server/main.go
\`\`\`

### 前端构建失败

\`\`\`bash
# 清除缓存重新安装
rm -rf node_modules package-lock.json
npm install
\`\`\`

### Docker容器无法启动

\`\`\`bash
# 查看容器日志
docker-compose logs backend
docker-compose logs postgres

# 重新构建
docker-compose build --no-cache
docker-compose up -d
\`\`\`

## 许可证

MIT License

## 联系方式

- 项目负责人: [Your Name]
- Email: [your.email@example.com]
- GitHub: [your-github-profile]

## 更新日志

### v1.0.0 (2026-01-16)
- ✨ 初始版本发布
- ✨ 完整的卡片管理功能
- ✨ 微信支付集成
- ✨ 管理后台
- ✨ Docker部署支持
