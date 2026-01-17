-- 插入测试SIM卡数据

-- 清空现有测试数据（可选）
-- DELETE FROM recharge_records;
-- DELETE FROM sim_cards WHERE card_no LIKE '898601234567890%';

-- 插入3张测试卡片
INSERT INTO sim_cards (card_no, device_no, start_date, expire_date, operator, package_type, remark)
VALUES
-- 卡片1：正常状态（距离到期还有1年）
('89860123456789012345', '866123456789012', '2026-01-01', '2027-01-01', '中国移动', '年卡', '测试卡片1 - 正常状态'),

-- 卡片2：即将到期（距离到期20天）
('89860123456789012346', '866123456789013', '2025-01-01', CURRENT_DATE + INTERVAL '20 days', '中国联通', '月卡', '测试卡片2 - 即将到期'),

-- 卡片3：已过期（已经过期30天）
('89860123456789012347', '866123456789014', '2024-01-01', CURRENT_DATE - INTERVAL '30 days', '中国电信', '年卡', '测试卡片3 - 已过期'),

-- 卡片4：正常状态（移动）
('89860123456789012348', '866123456789015', '2025-06-01', '2026-12-01', '中国移动', '年卡', '测试卡片4'),

-- 卡片5：正常状态（联通）
('89860123456789012349', '866123456789016', '2025-03-01', '2026-09-01', '中国联通', '年卡', '测试卡片5')
ON CONFLICT (card_no) DO NOTHING;

-- 插入一些充值记录（模拟已支付的订单）
INSERT INTO recharge_records (
    card_id,
    card_no,
    device_no,
    recharge_amount,
    recharge_years,
    old_expire_date,
    new_expire_date,
    trade_no,
    transaction_id,
    payment_status,
    paid_at,
    openid,
    ip_address
)
SELECT
    id,
    card_no,
    device_no,
    100.00,
    1,
    start_date,
    expire_date,
    'R2026011' || LPAD(id::text, 7, '0'),
    'wx2026011' || LPAD(id::text, 7, '0'),
    1,  -- 已支付
    CURRENT_TIMESTAMP - INTERVAL '5 days',
    'oTest_OpenID_' || id,
    '127.0.0.1'
FROM sim_cards
WHERE card_no IN ('89860123456789012345', '89860123456789012346')
ON CONFLICT (trade_no) DO NOTHING;

-- 查询插入结果
SELECT
    card_no,
    device_no,
    operator,
    expire_date,
    CASE
        WHEN status = 1 THEN '✅ 正常'
        WHEN status = 2 THEN '⚠️  即将到期'
        WHEN status = 3 THEN '❌ 已过期'
        ELSE '禁用'
    END as status_text,
    CASE
        WHEN expire_date >= CURRENT_DATE THEN (expire_date - CURRENT_DATE) || ' 天'
        ELSE '已过期 ' || (CURRENT_DATE - expire_date) || ' 天'
    END as remaining
FROM sim_cards
WHERE card_no LIKE '898601234567890%'
ORDER BY expire_date DESC;

-- 查询充值记录统计
SELECT
    COUNT(*) as total_recharge_count,
    COALESCE(SUM(recharge_amount), 0) as total_amount
FROM recharge_records
WHERE payment_status = 1;
