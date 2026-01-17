-- ========================================
-- IoT SIM卡管理系统 - 完整数据库初始化脚本
-- ========================================
-- 执行步骤:
-- 1. 打开PostgreSQL客户端: psql -U postgres
-- 2. 复制并执行以下所有SQL语句

-- ========================================
-- 第1步: 创建数据库和用户
-- ========================================

-- 创建数据库（如果不存在）
CREATE DATABASE iot_card_db;

-- 创建用户（如果不存在）
CREATE USER iot_user WITH PASSWORD 'adfhkIxcvYIK2189';

-- 授予权限
GRANT ALL PRIVILEGES ON DATABASE iot_card_db TO iot_user;

-- 切换到目标数据库
\c iot_card_db

-- 授予架构权限
GRANT ALL ON SCHEMA public TO iot_user;

-- 设置默认权限
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO iot_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO iot_user;

-- ========================================
-- 第2步: 创建管理员用户表
-- ========================================

CREATE TABLE IF NOT EXISTS admin_users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    real_name VARCHAR(100),
    phone VARCHAR(20),
    email VARCHAR(100),
    status SMALLINT DEFAULT 1, -- 1:启用 0:禁用
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_admin_username ON admin_users(username);
CREATE INDEX idx_admin_status ON admin_users(status);

-- 插入默认管理员账号 (密码: admin123)
-- bcrypt hash of "admin123"
INSERT INTO admin_users (username, password_hash, real_name, status) VALUES
('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', '系统管理员', 1);

COMMENT ON TABLE admin_users IS '管理员用户表';
COMMENT ON COLUMN admin_users.status IS '状态: 1=启用 0=禁用';

-- ========================================
-- 第3步: 创建SIM卡表
-- ========================================

CREATE TABLE IF NOT EXISTS sim_cards (
    id SERIAL PRIMARY KEY,
    card_no VARCHAR(50) UNIQUE NOT NULL,      -- 卡号(ICCID)
    device_no VARCHAR(50) UNIQUE NOT NULL,     -- 设备号(IMEI)
    start_date DATE NOT NULL,                  -- 开始日期
    expire_date DATE NOT NULL,                 -- 到期日期
    status SMALLINT DEFAULT 1,                 -- 1:正常 2:即将到期 3:已过期 0:禁用
    operator VARCHAR(20),                      -- 运营商(中国移动/联通/电信)
    package_type VARCHAR(50),                  -- 套餐类型
    total_recharge_count INT DEFAULT 0,        -- 总充值次数
    total_recharge_amount DECIMAL(10,2) DEFAULT 0, -- 总充值金额
    remark TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_card_no ON sim_cards(card_no);
CREATE INDEX idx_device_no ON sim_cards(device_no);
CREATE INDEX idx_expire_date ON sim_cards(expire_date);
CREATE INDEX idx_status ON sim_cards(status);

-- 自动更新状态的触发器
CREATE OR REPLACE FUNCTION update_sim_card_status()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.expire_date < CURRENT_DATE THEN
        NEW.status = 3; -- 已过期
    ELSIF NEW.expire_date <= CURRENT_DATE + INTERVAL '30 days' THEN
        NEW.status = 2; -- 即将到期
    ELSE
        NEW.status = 1; -- 正常
    END IF;
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_sim_card_status
BEFORE INSERT OR UPDATE ON sim_cards
FOR EACH ROW
EXECUTE FUNCTION update_sim_card_status();

COMMENT ON TABLE sim_cards IS 'SIM卡信息表';
COMMENT ON COLUMN sim_cards.status IS '状态: 1=正常 2=即将到期 3=已过期 0=禁用';
COMMENT ON COLUMN sim_cards.card_no IS '卡号(ICCID)';
COMMENT ON COLUMN sim_cards.device_no IS '设备号(IMEI)';

-- ========================================
-- 第4步: 创建充值记录表
-- ========================================

CREATE TABLE IF NOT EXISTS recharge_records (
    id SERIAL PRIMARY KEY,
    card_id INT NOT NULL,                      -- 关联sim_cards.id
    card_no VARCHAR(50) NOT NULL,              -- 冗余字段，便于查询
    device_no VARCHAR(50) NOT NULL,            -- 冗余字段
    recharge_amount DECIMAL(10,2) NOT NULL,    -- 充值金额
    recharge_years INT DEFAULT 1,              -- 充值年限
    old_expire_date DATE NOT NULL,             -- 原到期日期
    new_expire_date DATE NOT NULL,             -- 新到期日期
    payment_method VARCHAR(20) DEFAULT 'wechat', -- 支付方式
    trade_no VARCHAR(100) UNIQUE,              -- 交易单号(商户订单号)
    transaction_id VARCHAR(100),               -- 微信支付交易号
    payment_status SMALLINT DEFAULT 0,         -- 0:待支付 1:已支付 2:已退款 3:支付失败
    paid_at TIMESTAMP,                         -- 支付时间
    openid VARCHAR(100),                       -- 微信用户openid
    ip_address VARCHAR(50),                    -- 充值IP
    user_agent TEXT,                           -- 用户代理
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (card_id) REFERENCES sim_cards(id) ON DELETE CASCADE
);

CREATE INDEX idx_recharge_card_id ON recharge_records(card_id);
CREATE INDEX idx_recharge_card_no ON recharge_records(card_no);
CREATE INDEX idx_recharge_trade_no ON recharge_records(trade_no);
CREATE INDEX idx_recharge_payment_status ON recharge_records(payment_status);
CREATE INDEX idx_recharge_created_at ON recharge_records(created_at);
CREATE INDEX idx_recharge_paid_at ON recharge_records(paid_at);

-- 自动更新时间触发器
CREATE OR REPLACE FUNCTION update_recharge_record_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_recharge_record_timestamp
BEFORE UPDATE ON recharge_records
FOR EACH ROW
EXECUTE FUNCTION update_recharge_record_timestamp();

COMMENT ON TABLE recharge_records IS '充值记录表';
COMMENT ON COLUMN recharge_records.payment_status IS '支付状态: 0=待支付 1=已支付 2=已退款 3=支付失败';

-- ========================================
-- 第5步: 创建系统配置表
-- ========================================

CREATE TABLE IF NOT EXISTS system_config (
    id SERIAL PRIMARY KEY,
    config_key VARCHAR(50) UNIQUE NOT NULL,
    config_value TEXT NOT NULL,
    config_type VARCHAR(20) DEFAULT 'string', -- string/number/json
    description VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_config_key ON system_config(config_key);

-- 自动更新时间触发器
CREATE OR REPLACE FUNCTION update_system_config_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_system_config_timestamp
BEFORE UPDATE ON system_config
FOR EACH ROW
EXECUTE FUNCTION update_system_config_timestamp();

-- 插入默认配置
INSERT INTO system_config (config_key, config_value, config_type, description) VALUES
('recharge_price', '100.00', 'number', '充值价格(元/年)'),
('wechat_appid', '', 'string', '微信公众号AppID'),
('wechat_mch_id', '', 'string', '微信商户号'),
('wechat_api_key', '', 'string', '微信API密钥'),
('wechat_api_v3_key', '', 'string', '微信APIv3密钥'),
('wechat_serial_no', '', 'string', '微信证书序列号'),
('notify_url', 'https://yourdomain.com/api/v1/payment/notify', 'string', '微信支付回调地址'),
('alert_days', '30', 'number', '到期提醒天数');

-- ========================================
-- 第6步: 创建统计视图（可选）
-- ========================================

CREATE OR REPLACE VIEW v_statistics AS
SELECT
    (SELECT COUNT(*) FROM sim_cards) as total_cards,
    (SELECT COUNT(*) FROM sim_cards WHERE status = 1) as active_cards,
    (SELECT COUNT(*) FROM sim_cards WHERE status = 2) as expiring_cards,
    (SELECT COUNT(*) FROM sim_cards WHERE status = 3) as expired_cards,
    (SELECT COALESCE(SUM(recharge_amount), 0) FROM recharge_records WHERE payment_status = 1) as total_recharge_amount,
    (SELECT COUNT(*) FROM recharge_records WHERE payment_status = 1) as total_recharge_count,
    (SELECT COALESCE(SUM(recharge_amount), 0) FROM recharge_records 
     WHERE payment_status = 1 AND EXTRACT(YEAR FROM paid_at) = EXTRACT(YEAR FROM CURRENT_DATE)) as current_year_amount;

-- ========================================
-- 完成！
-- ========================================
-- 数据库信息:
--   数据库名: iot_card_db
--   用户名: iot_user
--   密码: adfhkIxcvYIK2189
--
-- 已创建的表:
--   - admin_users (管理员用户表)
--   - sim_cards (SIM卡信息表)
--   - recharge_records (充值记录表)
--   - system_config (系统配置表)
--
-- 默认管理员账号:
--   用户名: admin
--   密码: admin123
-- ========================================
