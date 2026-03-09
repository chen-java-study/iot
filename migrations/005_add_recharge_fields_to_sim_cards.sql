-- 005_add_recharge_fields_to_sim_cards.sql
-- 为 sim_cards 表添加最近充值时间和充值金额字段

ALTER TABLE sim_cards ADD COLUMN IF NOT EXISTS last_recharge_time TIMESTAMP;
ALTER TABLE sim_cards ADD COLUMN IF NOT EXISTS last_recharge_amount DECIMAL(10,2) DEFAULT 0;

COMMENT ON COLUMN sim_cards.last_recharge_time IS '最近充值时间';
COMMENT ON COLUMN sim_cards.last_recharge_amount IS '最近充值金额';
