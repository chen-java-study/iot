-- 004_create_system_config.sql
-- 系统配置表

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

-- 创建统计视图
CREATE OR REPLACE VIEW v_statistics AS
SELECT
    (SELECT COUNT(*) FROM sim_cards) as total_cards,
    (SELECT COUNT(*) FROM sim_cards WHERE status = 1) as active_cards,
    (SELECT COUNT(*) FROM sim_cards WHERE status = 2) as expiring_cards,
    (SELECT COUNT(*) FROM sim_cards WHERE status = 3) as expired_cards,
    (SELECT COALESCE(SUM(recharge_amount), 0) FROM recharge_records WHERE payment_status = 1) as total_recharge_amount,
    (SELECT COUNT(*) FROM recharge_records WHERE payment_status = 1) as total_recharge_count,
    (SELECT COALESCE(SUM(recharge_amount), 0) FROM recharge_records
     WHERE payment_status = 1 AND DATE(paid_at) = CURRENT_DATE) as today_amount,
    (SELECT COALESCE(SUM(recharge_amount), 0) FROM recharge_records
     WHERE payment_status = 1 AND DATE_TRUNC('month', paid_at) = DATE_TRUNC('month', CURRENT_DATE)) as month_amount;

COMMENT ON TABLE system_config IS '系统配置表';
COMMENT ON VIEW v_statistics IS '统计数据视图';
