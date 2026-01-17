# 数据库初始化指南

## 概述

本指南说明如何为 IoT SIM卡管理系统创建数据库和表。系统使用PostgreSQL数据库，包含以下4个主要表：
- **admin_users**: 管理员用户信息
- **sim_cards**: SIM卡信息
- **recharge_records**: 充值记录
- **system_config**: 系统配置

---

## 前置条件

1. **PostgreSQL已安装并运行**
   - Windows: PostgreSQL 12+
   - Linux/Mac: PostgreSQL 12+
   - 默认端口: 5432

2. **psql命令行工具可用**
   ```bash
   # 验证安装
   psql --version
   ```

---

## 方式A: 自动初始化（推荐）

### Windows系统

1. 打开命令提示符或PowerShell
2. 进入项目目录
3. 运行初始化脚本：
   ```bash
   cd backend\migrations
   init_db_windows.bat
   ```

### Linux/Mac系统

1. 打开终端
2. 进入项目目录
3. 运行初始化脚本：
   ```bash
   cd backend/migrations
   chmod +x init_db.sh
   ./init_db.sh
   ```

---

## 方式B: 手动逐步初始化

### 步骤1: 创建数据库和用户

```bash
# 连接到PostgreSQL默认数据库
psql -U postgres

# 在psql中执行以下命令
CREATE DATABASE iot_card_db;
CREATE USER iot_user WITH PASSWORD 'adfhkIxcvYIK2189';
GRANT ALL PRIVILEGES ON DATABASE iot_card_db TO iot_user;
ALTER DATABASE iot_card_db OWNER TO iot_user;
\q  # 退出psql
```

### 步骤2: 执行迁移脚本

按顺序执行SQL迁移脚本：

```bash
cd backend/migrations

# 创建管理员用户表
psql -U iot_user -h localhost -d iot_card_db -f 001_create_admin_users.sql

# 创建SIM卡表
psql -U iot_user -h localhost -d iot_card_db -f 002_create_sim_cards.sql

# 创建充值记录表
psql -U iot_user -h localhost -d iot_card_db -f 003_create_recharge_records.sql

# 创建系统配置表
psql -U iot_user -h localhost -d iot_card_db -f 004_create_system_config.sql

# 插入测试数据（可选）
psql -U iot_user -h localhost -d iot_card_db -f insert_test_data.sql
```

---

## 方式C: 使用Docker

如果使用Docker运行PostgreSQL：

```bash
# 1. 启动PostgreSQL容器
docker run -d \
  --name iot_postgres \
  -e POSTGRES_DB=iot_card_db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres123 \
  -p 5432:5432 \
  postgres:15-alpine

# 2. 等待容器启动（约10秒）
sleep 10

# 3. 执行迁移脚本
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend/migrations/001_create_admin_users.sql
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend/migrations/002_create_sim_cards.sql
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend/migrations/003_create_recharge_records.sql
docker exec -i iot_postgres psql -U postgres -d iot_card_db < backend/migrations/004_create_system_config.sql
```

---

## 验证数据库初始化

### 验证表是否成功创建

```bash
# 连接到数据库
psql -U iot_user -h localhost -d iot_card_db

# 在psql中查看所有表
\dt

# 查看表结构
\d admin_users
\d sim_cards
\d recharge_records
\d system_config

# 查看数据
SELECT * FROM admin_users;
SELECT * FROM system_config;

# 退出
\q
```

### 预期输出

表格如下所示应该存在：
- `admin_users` - 默认包含1条admin用户记录
- `sim_cards` - 可能包含测试数据
- `recharge_records` - 充值记录表
- `system_config` - 包含系统配置

---

## 数据库凭证

使用以下凭证连接到数据库：

| 项目 | 值 |
|------|-----|
| 主机 | localhost |
| 端口 | 5432 |
| 数据库 | iot_card_db |
| 用户 | iot_user |
| 密码 | adfhkIxcvYIK2189 |

**注意**: 如果使用Docker，请根据Docker配置调整凭证。

---

## 更新backend配置

确保 `backend/configs/config.yaml` 中的数据库配置与上述凭证匹配：

```yaml
database:
  host: 127.0.0.1
  port: 5432
  user: iot_user
  password: "adfhkIxcvYIK2189"
  dbname: iot_card_db
```

---

## 常见问题解决

### 问题1: psql命令不找到
**解决方案**: 
- Windows: 将PostgreSQL的bin目录添加到系统PATH环境变量
- Linux/Mac: 确保PostgreSQL已正确安装

### 问题2: 连接被拒绝
**解决方案**:
- 检查PostgreSQL服务是否运行
- 检查主机名、端口、用户名和密码
- 确保防火墙允许5432端口

### 问题3: 数据库已存在
**解决方案**:
```bash
# 删除现有数据库（谨慎操作）
psql -U postgres -c "DROP DATABASE iot_card_db;"

# 重新运行初始化脚本
```

### 问题4: 权限错误
**解决方案**:
```bash
# 检查iot_user是否有权限
psql -U postgres

# 授予权限
ALTER USER iot_user CREATEDB;
GRANT ALL PRIVILEGES ON DATABASE iot_card_db TO iot_user;
```

---

## 表结构说明

### admin_users 表
- 存储管理员账户信息
- 默认包含一个admin用户（密码: admin123）
- 使用bcrypt加密密码

### sim_cards 表
- 存储SIM卡信息
- 包含卡号(ICCID)和设备号(IMEI)
- 自动更新卡片过期状态的触发器

### recharge_records 表
- 存储充值交易记录
- 包含微信支付信息
- 与sim_cards表有外键关联

### system_config 表
- 存储系统配置信息
- 包含微信支付配置、充值价格等
- 默认包含必要的配置项

---

## 下一步

初始化完成后，可以：
1. 启动后端服务: `go run cmd/server/main.go`
2. 启动前端开发服务器
3. 开始功能测试

详见项目的 TESTING_GUIDE.md 文件。
