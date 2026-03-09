-- 003_create_recharge_records.sql
-- 充值记录表

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
