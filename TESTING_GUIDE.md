# 前端测试指南

## 快速测试方案（推荐）

由于完整的微信支付需要配置，我们提供**模拟数据测试**方案，可以快速验证所有功能。

---

## 方案一：前后端联调测试（完整功能）

### 步骤1：启动PostgreSQL数据库

**方式A - 使用Docker（推荐）：**
```bash
# 启动PostgreSQL容器
docker run -d \
  --name iot_postgres \
  -e POSTGRES_DB=iot_card_db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres123 \
  -p 5432:5432 \
  postgres:15-alpine

# 等待数据库启动（约10秒）
sleep 10

# 执行数据库迁移
docker cp backend/migrations iot_postgres:/tmp/
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend/migrations/001_create_admin_users.sql
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend/migrations/002_create_sim_cards.sql
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend/migrations/003_create_recharge_records.sql
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend/migrations/004_create_system_config.sql
```

**方式B - 使用本地PostgreSQL：**
```bash
# 确保PostgreSQL已安装并启动
# 创建数据库
createdb iot_card_db

# 执行迁移脚本
cd backend/migrations
psql -d iot_card_db -f 001_create_admin_users.sql
psql -d iot_card_db -f 002_create_sim_cards.sql
psql -d iot_card_db -f 003_create_recharge_records.sql
psql -d iot_card_db -f 004_create_system_config.sql
cd ../..
```

### 步骤2：配置后端

编辑 `backend/configs/config.yaml`：
```yaml
server:
  port: 8080
  mode: debug  # 开发模式，显示详细日志

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres123  # 改为你的密码
  dbname: iot_card_db

jwt:
  secret_key: test_secret_key_for_development
  expire_hours: 24

wechat:
  app_id: ""  # 暂时留空，测试时使用模拟数据
  mch_id: ""
  api_v3_key: ""
  serial_no: ""
  private_key_path: ""
  notify_url: "http://localhost:8080/api/v1/payment/notify"
```

### 步骤3：启动后端服务

**Windows (PowerShell/CMD):**
```bash
cd backend
go run cmd/server/main.go
```

**Linux/Mac:**
```bash
cd backend
go run cmd/server/main.go
```

**看到以下输出表示成功：**
```
数据库连接成功
服务器启动在端口 8080
[GIN-debug] Listening and serving HTTP on :8080
```

### 步骤4：启动H5前端

**新开一个终端窗口：**
```bash
cd frontend/h5
npm install
npm run dev
```

**成功启动后会显示：**
```
VITE v5.x.x  ready in xxx ms

➜  Local:   http://localhost:3000/
➜  Network: use --host to expose
```

### 步骤5：启动管理端

**再新开一个终端窗口：**
```bash
cd frontend/admin
npm install
npm run dev
```

**成功启动后会显示：**
```
VITE v5.x.x  ready in xxx ms

➜  Local:   http://localhost:3001/
➜  Network: use --host to expose
```

---

## 测试流程

### 一、测试管理端 (http://localhost:3001)

#### 1. 登录功能测试
1. 浏览器打开 `http://localhost:3001`
2. 会自动跳转到登录页 `http://localhost:3001/login`
3. 输入默认账号：
   - 用户名：`admin`
   - 密码：`admin123`
4. 点击"登录"按钮
5. **预期结果**：
   - 提示"登录成功"
   - 跳转到首页仪表盘
   - 显示统计数据（初始全为0）

**调试技巧：**
- 打开浏览器开发者工具 (F12)
- 查看 Network 标签，确认请求发送到 `http://localhost:8080/api/v1/admin/login`
- 如果报错，查看 Console 标签的错误信息

#### 2. 卡片管理测试

**添加测试卡片：**
1. 点击左侧菜单"卡片管理"
2. 点击"添加卡片"按钮（需要手动实现添加功能，或使用API）

**使用API直接添加测试数据：**

**方式A - 使用curl（推荐）：**
```bash
# 先登录获取token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

# 添加测试卡片
curl -X POST http://localhost:8080/api/v1/admin/cards \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "card_no": "89860123456789012345",
    "device_no": "866123456789012",
    "start_date": "2026-01-01T00:00:00Z",
    "expire_date": "2027-01-01T00:00:00Z",
    "operator": "中国移动",
    "package_type": "年卡",
    "remark": "测试卡片"
  }'
```

**方式B - 使用Postman/Insomnia：**
1. POST `http://localhost:8080/api/v1/admin/login`
   ```json
   {
     "username": "admin",
     "password": "admin123"
   }
   ```
2. 复制返回的token
3. POST `http://localhost:8080/api/v1/admin/cards`
   - Headers: `Authorization: Bearer <你的token>`
   ```json
   {
     "card_no": "89860123456789012345",
     "device_no": "866123456789012",
     "start_date": "2026-01-01T00:00:00Z",
     "expire_date": "2027-01-01T00:00:00Z",
     "operator": "中国移动",
     "package_type": "年卡"
   }
   ```

**方式C - 直接插入数据库：**
```bash
docker exec -it iot_postgres psql -U postgres -d iot_card_db
# 或本地 psql -d iot_card_db

# 执行SQL
INSERT INTO sim_cards (card_no, device_no, start_date, expire_date, operator, package_type)
VALUES
('89860123456789012345', '866123456789012', '2026-01-01', '2027-01-01', '中国移动', '年卡'),
('89860123456789012346', '866123456789013', '2025-12-01', '2026-01-20', '中国联通', '月卡'),
('89860123456789012347', '866123456789014', '2025-01-01', '2025-12-31', '中国电信', '年卡');
```

3. 刷新管理端页面，应该能看到卡片列表
4. **预期结果**：
   - 显示卡片列表
   - 可以看到卡号、设备号、到期日期、状态等信息

#### 3. 首页统计测试
1. 点击左侧菜单"首页仪表盘"
2. **预期结果**：
   - 卡片总数：显示已添加的卡片数量
   - 充值总金额：0（还没有充值记录）
   - 今日收入：0
   - 本月收入：0

---

### 二、测试H5端 (http://localhost:3000)

#### 1. 查询功能测试

1. 浏览器打开 `http://localhost:3000`
2. 进入查询页面
3. 在搜索框输入：`89860123456789012345`（你添加的卡号）
4. 点击"查询"按钮
5. **预期结果**：
   - 显示卡片信息卡片
   - 显示：卡号、设备号、运营商、到期时间、剩余天数、状态
   - 状态应该显示为"正常"（绿色）或"即将到期"（橙色）

**测试不同状态的卡：**
- 正常卡：到期日期 > 今天+30天
- 即将到期：到期日期在今天到今天+30天之间
- 已过期：到期日期 < 今天

#### 2. 充值功能测试

1. 在查询结果页面，点击"充值续费"按钮
2. 跳转到充值页面
3. 显示：卡号、当前到期时间、充值金额（¥100.00）
4. 点击"立即支付"按钮
5. **预期结果**（当前是模拟模式）：
   - 显示"订单创建成功，订单号：Rxxxxxxxx"
   - 实际项目中这里会调起微信支付

**查看充值记录：**
```bash
# 查看数据库中的充值记录
docker exec -it iot_postgres psql -U postgres -d iot_card_db -c "SELECT * FROM recharge_records;"
```

---

## 方案二：纯前端测试（不连后端）

如果后端配置有问题，可以先测试前端界面和交互：

### 修改前端API为Mock模式

**H5端 - 修改 `frontend/h5/src/api/card.js`：**
```javascript
export function queryCard(keyword) {
  // 返回模拟数据
  return Promise.resolve({
    id: 1,
    card_no: '89860123456789012345',
    device_no: '866123456789012',
    start_date: '2026-01-01',
    expire_date: '2027-01-01',
    status: 1,
    status_text: '正常',
    operator: '中国移动',
    days_remaining: 350
  })
}
```

**管理端 - 修改 `frontend/admin/src/views/Login.vue`：**
```javascript
const handleLogin = async () => {
  // 模拟登录成功
  localStorage.setItem('admin_token', 'mock_token_12345')
  localStorage.setItem('admin_user', JSON.stringify({
    id: 1,
    username: 'admin',
    real_name: '管理员'
  }))
  ElMessage.success('登录成功')
  router.push('/')
}
```

然后正常启动前端进行界面测试。

---

## 常见问题解决

### 问题1: 后端启动失败 - "连接数据库失败"

**解决方法：**
```bash
# 检查PostgreSQL是否运行
docker ps | grep postgres
# 或
psql -U postgres -l

# 检查配置文件中的数据库密码是否正确
cat backend/configs/config.yaml

# 测试数据库连接
psql -h localhost -U postgres -d iot_card_db
```

### 问题2: 前端请求失败 - "Network Error"

**检查后端是否运行：**
```bash
# 访问后端健康检查
curl http://localhost:8080/api/v1/card/query?keyword=test
```

**检查CORS配置：**
打开浏览器Console，如果看到CORS错误，确认后端已启用CORS中间件（已在router.go中配置）。

### 问题3: 前端npm install失败

**解决方法：**
```bash
# 清除缓存
npm cache clean --force
rm -rf node_modules package-lock.json

# 重新安装
npm install

# 或使用淘宝镜像
npm install --registry=https://registry.npmmirror.com
```

### 问题4: Go依赖下载慢

**配置国内代理：**
```bash
go env -w GOPROXY=https://goproxy.cn,direct
go mod download
```

---

## 完整测试检查清单

### 管理端测试清单
- [ ] 登录页面显示正常
- [ ] 能够成功登录（admin/admin123）
- [ ] 未登录访问其他页面会跳转到登录页
- [ ] 首页仪表盘显示统计数据
- [ ] 卡片列表页显示数据
- [ ] 充值记录页显示数据
- [ ] 系统配置页显示配置项
- [ ] 退出登录功能正常

### H5端测试清单
- [ ] 查询页面显示正常
- [ ] 输入卡号能够查询到数据
- [ ] 卡片信息显示完整
- [ ] 状态显示正确（正常/即将到期/已过期）
- [ ] 充值按钮点击正常
- [ ] 充值页面显示正确信息

### API测试清单
- [ ] POST /api/v1/admin/login - 登录成功
- [ ] GET /api/v1/admin/statistics - 获取统计
- [ ] GET /api/v1/admin/cards - 获取卡片列表
- [ ] POST /api/v1/admin/cards - 创建卡片
- [ ] GET /api/v1/card/query - 查询卡片
- [ ] POST /api/v1/payment/create - 创建订单

---

## 测试演示视频脚本

如果需要录制演示，可以按以下流程：

1. **管理端演示（2分钟）**
   - 登录系统
   - 查看首页统计
   - 添加一张测试卡
   - 查看卡片列表

2. **H5端演示（1分钟）**
   - 打开查询页面
   - 输入卡号查询
   - 查看卡片信息
   - 点击充值按钮

3. **数据验证（30秒）**
   - 查看数据库记录
   - 查看后端日志

---

## 下一步：接入真实微信支付

测试通过后，如需接入真实微信支付：

1. 在微信公众平台配置支付授权目录
2. 配置 `backend/configs/config.yaml` 中的微信参数
3. 上传API证书文件
4. 修改 `frontend/h5/src/views/Recharge.vue` 调用真实的WeixinJSBridge
5. 部署到HTTPS域名（微信支付要求）

需要我帮你配置吗？
