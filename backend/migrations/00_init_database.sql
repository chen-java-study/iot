-- 00_init_database.sql
-- 初始化数据库脚本（Windows/Linux通用）
-- 使用方法：psql -U postgres -f 00_init_database.sql

-- 创建数据库（如果不存在）
SELECT 'CREATE DATABASE iot_card_db' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'iot_card_db')\gexec

-- 创建用户（如果不存在）
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_user WHERE usename = 'iot_user') THEN
        CREATE USER iot_user WITH PASSWORD 'adfhkIxcvYIK2189';
    END IF;
END
$$;

-- 授予权限
GRANT ALL PRIVILEGES ON DATABASE iot_card_db TO iot_user;

\c iot_card_db

-- 授予架构权限
GRANT ALL ON SCHEMA public TO iot_user;

-- 设置默认权限
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO iot_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO iot_user;
