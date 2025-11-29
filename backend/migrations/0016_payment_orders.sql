-- Adds payment_orders table for recharge via 易支付
CREATE TABLE IF NOT EXISTS payment_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT NOT NULL,
    out_trade_no TEXT UNIQUE NOT NULL,
    provider_trade_no TEXT DEFAULT '',
    pay_type TEXT DEFAULT 'alipay',
    status TEXT NOT NULL DEFAULT 'pending', -- pending | paid | failed
    money_cents BIGINT NOT NULL,
    coins BIGINT NOT NULL,
    notify_payload JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_payment_orders_user ON payment_orders(user_id);

COMMENT ON TABLE payment_orders IS '充值订单';
COMMENT ON COLUMN payment_orders.money_cents IS '支付金额，单位分';
COMMENT ON COLUMN payment_orders.coins IS '到账虚拟币数量';

