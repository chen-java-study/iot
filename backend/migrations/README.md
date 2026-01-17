# 数据库迁移文件说明

本目录包含所有数据库初始化和迁移脚本。

## 文件说明

| 文件 | 说明 |
|------|------|
| **00_init_database.sql** | 初始化数据库和用户的SQL脚本 |
| **001_create_admin_users.sql** | 创建管理员用户表 |
| **002_create_sim_cards.sql** | 创建SIM卡信息表 |
| **003_create_recharge_records.sql** | 创建充值记录表 |
| **004_create_system_config.sql** | 创建系统配置表 |
| **insert_test_data.sql** | 插入测试数据 |
| **init_db_windows.bat** | Windows自动初始化脚本 |
| **init_db.sh** | Linux/Mac自动初始化脚本 |
| **DATABASE_SETUP.md** | 详细的数据库设置指南 |

## 快速开始

### Windows用户
```bash
cd backend\migrations
init_db_windows.bat
```

### Linux/Mac用户
```bash
cd backend/migrations
chmod +x init_db.sh
./init_db.sh
```

## 详细说明

请参考 **DATABASE_SETUP.md** 获取完整的设置指南，包括：
- 前置条件检查
- 多种初始化方式
- 验证方法
- 常见问题解决

## 手动执行

如需手动执行迁移脚本：

```bash
# 1. 先执行数据库初始化
psql -U postgres -f 00_init_database.sql

# 2. 然后按顺序执行表创建脚本
psql -U iot_user -d iot_card_db -f 001_create_admin_users.sql
psql -U iot_user -d iot_card_db -f 002_create_sim_cards.sql
psql -U iot_user -d iot_card_db -f 003_create_recharge_records.sql
psql -U iot_user -d iot_card_db -f 004_create_system_config.sql

# 3. 可选：插入测试数据
psql -U iot_user -d iot_card_db -f insert_test_data.sql
```

## 数据库凭证

- **数据库**: iot_card_db
- **用户**: iot_user
- **密码**: adfhkIxcvYIK2189
- **主机**: localhost
- **端口**: 5432

确保 `backend/configs/config.yaml` 中的配置与此相同。
