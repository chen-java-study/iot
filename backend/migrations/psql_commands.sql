-- ========================================
-- PostgreSQL 数据库查询命令
-- ========================================

-- 1. 连接到PostgreSQL（在命令行中执行）
psql -U postgres

-- 2. 在psql提示符下，切换到目标数据库
\c iot_card_db

-- 3. 查看所有表
\dt

-- 4. 查看表详细结构（按需选择）

-- 查看admin_users表结构
\d admin_users

-- 查看sim_cards表结构
\d sim_cards

-- 查看recharge_records表结构
\d recharge_records

-- 查看system_config表结构
\d system_config

-- ========================================
-- 5. 查询表中的数据
-- ========================================

-- 查询管理员用户
SELECT * FROM admin_users;

-- 查询SIM卡信息
SELECT * FROM sim_cards;

-- 查询充值记录
SELECT * FROM recharge_records;

-- 查询系统配置
SELECT * FROM system_config;

-- ========================================
-- 6. 其他常用查询
-- ========================================

-- 查询表数量统计
SELECT COUNT(*) as total_cards FROM sim_cards;
SELECT COUNT(*) as total_recharge FROM recharge_records;
SELECT COUNT(*) as total_config FROM system_config;

-- 查询统计视图
SELECT * FROM v_statistics;

-- 查看所有索引
SELECT indexname FROM pg_indexes WHERE tablename = 'sim_cards';

-- 查看用户列表
\du

-- 查看数据库列表
\l

-- 退出psql
\q
