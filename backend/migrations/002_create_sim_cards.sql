-- 002_create_sim_cards.sql
-- SIM卡表

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
