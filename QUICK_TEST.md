# 快速测试命令

## 一、准备工作

### Windows用户

#### 方式1：双击启动（最简单）
```
双击运行：start-dev.bat
```

#### 方式2：Docker快速启动数据库
```bash
# 启动PostgreSQL（首次运行）
docker run -d --name iot_postgres ^
  -e POSTGRES_DB=iot_card_db ^
  -e POSTGRES_USER=postgres ^
  -e POSTGRES_PASSWORD=postgres123 ^
  -p 5432:5432 ^
  postgres:15-alpine

# 等待10秒让数据库启动
timeout /t 10

# 执行数据库迁移
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend/migrations/001_create_admin_users.sql
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend/migrations/002_create_sim_cards.sql
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend/migrations/003_create_recharge_records.sql
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend/migrations/004_create_system_config.sql

# 插入测试数据
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend/migrations/insert_test_data.sql
```

### Linux/Mac用户
```bash
# 一键启动
chmod +x start-dev.sh
./start-dev.sh
```

---

## 二、插入测试数据

### 方式1：使用SQL文件（推荐）
```bash
# Docker数据库
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend/migrations/insert_test_data.sql

# 本地数据库
psql -U postgres -d iot_card_db -f backend/migrations/insert_test_data.sql
```

### 方式2：使用curl命令

#### Step 1: 登录获取Token
```bash
# Windows (PowerShell)
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/admin/login" -Method Post -ContentType "application/json" -Body '{"username":"admin","password":"admin123"}'
$token = $response.data.token
echo $token

# Linux/Mac
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | jq -r '.data.token')
echo $TOKEN
```

#### Step 2: 添加测试卡片
```bash
# Windows (PowerShell)
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/admin/cards" `
  -Method Post `
  -Headers @{"Authorization"="Bearer $token"; "Content-Type"="application/json"} `
  -Body '{
    "card_no": "89860123456789012345",
    "device_no": "866123456789012",
    "start_date": "2026-01-01T00:00:00Z",
    "expire_date": "2027-01-01T00:00:00Z",
    "operator": "中国移动",
    "package_type": "年卡",
    "remark": "PowerShell测试卡片"
  }'

# Linux/Mac
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
    "remark": "curl测试卡片"
  }'
```

---

## 三、API测试命令

### 1. 查询卡片（无需认证）
```bash
# Windows
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/card/query?keyword=89860123456789012345"

# Linux/Mac
curl "http://localhost:8080/api/v1/card/query?keyword=89860123456789012345"
```

### 2. 管理员登录
```bash
# Windows
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/admin/login" `
  -Method Post `
  -ContentType "application/json" `
  -Body '{"username":"admin","password":"admin123"}'

# Linux/Mac
curl -X POST http://localhost:8080/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 3. 获取统计数据（需要Token）
```bash
# Windows
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/admin/statistics" `
  -Headers @{"Authorization"="Bearer $token"}

# Linux/Mac
curl http://localhost:8080/api/v1/admin/statistics \
  -H "Authorization: Bearer $TOKEN"
```

### 4. 获取卡片列表
```bash
# Windows
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/admin/cards?page=1&page_size=10" `
  -Headers @{"Authorization"="Bearer $token"}

# Linux/Mac
curl "http://localhost:8080/api/v1/admin/cards?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"
```

### 5. 创建充值订单
```bash
# Windows
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/payment/create" `
  -Method Post `
  -ContentType "application/json" `
  -Body '{"card_no":"89860123456789012345","openid":"test_openid_123"}'

# Linux/Mac
curl -X POST http://localhost:8080/api/v1/payment/create \
  -H "Content-Type: application/json" \
  -d '{"card_no":"89860123456789012345","openid":"test_openid_123"}'
```

---

## 四、浏览器测试

### H5端测试流程
1. 打开浏览器访问：`http://localhost:3000`
2. 输入测试卡号：`89860123456789012345`
3. 点击"查询"
4. 查看卡片信息
5. 点击"充值续费"
6. 查看充值页面

### 管理端测试流程
1. 打开浏览器访问：`http://localhost:3001`
2. 会自动跳转到登录页
3. 输入账号密码：`admin` / `admin123`
4. 点击"登录"
5. 进入首页查看统计数据
6. 点击"卡片管理"查看卡片列表
7. 点击"充值记录"查看充值数据

---

## 五、数据库直接查询

### 查看所有卡片
```bash
# Docker
docker exec -it iot_postgres psql -U postgres -d iot_card_db -c "SELECT card_no, device_no, operator, expire_date, status FROM sim_cards;"

# 本地
psql -U postgres -d iot_card_db -c "SELECT card_no, device_no, operator, expire_date, status FROM sim_cards;"
```

### 查看充值记录
```bash
docker exec -it iot_postgres psql -U postgres -d iot_card_db -c "SELECT trade_no, card_no, recharge_amount, payment_status, created_at FROM recharge_records ORDER BY created_at DESC LIMIT 10;"
```

### 查看统计数据
```bash
docker exec -it iot_postgres psql -U postgres -d iot_card_db -c "SELECT * FROM v_statistics;"
```

---

## 六、调试技巧

### 查看后端日志
后端日志会直接输出到终端，查看启动后端的窗口即可。

### 查看前端日志
打开浏览器开发者工具（F12）：
- **Console**: 查看JavaScript错误
- **Network**: 查看API请求和响应
- **Application**: 查看LocalStorage中的Token

### 常用调试命令

#### 测试后端是否运行
```bash
# Windows
Invoke-WebRequest -Uri "http://localhost:8080/api/v1/card/query?keyword=test"

# Linux/Mac
curl http://localhost:8080/api/v1/card/query?keyword=test
```

#### 查看端口占用
```bash
# Windows
netstat -ano | findstr :8080
netstat -ano | findstr :3000
netstat -ano | findstr :3001

# Linux/Mac
lsof -i :8080
lsof -i :3000
lsof -i :3001
```

#### 清理重启
```bash
# 停止所有服务
# Windows: 双击 stop-dev.bat
# Linux/Mac: Ctrl+C

# 清除前端缓存
cd frontend/h5
rm -rf node_modules package-lock.json
npm install

cd ../admin
rm -rf node_modules package-lock.json
npm install
```

---

## 七、测试检查清单

### 后端测试
- [ ] Go服务启动成功（端口8080）
- [ ] 数据库连接成功
- [ ] API可以访问
- [ ] 登录接口返回Token
- [ ] 查询接口返回数据

### H5前端测试
- [ ] 页面正常显示（localhost:3000）
- [ ] 查询功能正常
- [ ] 显示卡片信息
- [ ] 充值页面显示正确
- [ ] 无控制台错误

### 管理端测试
- [ ] 页面正常显示（localhost:3001）
- [ ] 登录功能正常
- [ ] 路由守卫生效（未登录跳转）
- [ ] 首页统计显示
- [ ] 卡片列表显示
- [ ] 无控制台错误

---

## 八、完整测试脚本（复制粘贴运行）

### Windows PowerShell完整测试
```powershell
# 1. 启动数据库
docker run -d --name iot_postgres -e POSTGRES_DB=iot_card_db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres123 -p 5432:5432 postgres:15-alpine
Start-Sleep -Seconds 10

# 2. 执行迁移
Get-ChildItem backend\migrations\*.sql | ForEach-Object {
    Get-Content $_.FullName | docker exec -i iot_postgres psql -U postgres -d iot_card_db
}

# 3. 测试API
Write-Host "测试登录..." -ForegroundColor Green
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/admin/login" -Method Post -ContentType "application/json" -Body '{"username":"admin","password":"admin123"}'
$token = $response.data.token
Write-Host "Token: $token" -ForegroundColor Yellow

Write-Host "`n测试统计..." -ForegroundColor Green
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/admin/statistics" -Headers @{"Authorization"="Bearer $token"}

Write-Host "`n测试查询..." -ForegroundColor Green
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/card/query?keyword=89860123456789012345"

Write-Host "`n✅ 测试完成！" -ForegroundColor Green
```

### Linux/Mac完整测试
```bash
#!/bin/bash

# 1. 启动数据库
docker run -d --name iot_postgres \
  -e POSTGRES_DB=iot_card_db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres123 \
  -p 5432:5432 \
  postgres:15-alpine
sleep 10

# 2. 执行迁移
for sql in backend/migrations/*.sql; do
    docker exec -i iot_postgres psql -U postgres -d iot_card_db < "$sql"
done

# 3. 测试API
echo "测试登录..."
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | jq -r '.data.token')
echo "Token: $TOKEN"

echo -e "\n测试统计..."
curl -s http://localhost:8080/api/v1/admin/statistics \
  -H "Authorization: Bearer $TOKEN" | jq

echo -e "\n测试查询..."
curl -s "http://localhost:8080/api/v1/card/query?keyword=89860123456789012345" | jq

echo -e "\n✅ 测试完成！"
```

---

## 需要帮助？

如遇到问题，请查看：
1. `TESTING_GUIDE.md` - 详细测试指南
2. `README.md` - 完整项目文档
3. 后端日志输出
4. 浏览器控制台错误信息

常见问题：
- 端口被占用：修改配置文件中的端口号
- 数据库连接失败：检查PostgreSQL是否运行
- 前端请求失败：检查后端是否启动
- CORS错误：已在后端配置，清除浏览器缓存重试
