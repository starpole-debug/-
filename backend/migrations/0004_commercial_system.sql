-- +goose Up
-- VIP Tiers Configuration
CREATE TABLE IF NOT EXISTS vip_tiers (
    id TEXT PRIMARY KEY, -- e.g., 'vip_monthly', 'vip_yearly'
    name TEXT NOT NULL,
    price INTEGER NOT NULL, -- in cents
    monthly_tickets INTEGER NOT NULL DEFAULT 0,
    tts_quota_daily INTEGER NOT NULL DEFAULT 0,
    img_gen_quota_daily INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- User Assets (Wallet)
CREATE TABLE IF NOT EXISTS user_assets (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    balance INTEGER NOT NULL DEFAULT 0, -- Tokens
    monthly_tickets INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- User Quotas (Daily/Monthly Limits)
CREATE TABLE IF NOT EXISTS user_quotas (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    tts_remaining INTEGER NOT NULL DEFAULT 0,
    img_gen_remaining INTEGER NOT NULL DEFAULT 0,
    next_reset_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- VIP Subscriptions
CREATE TABLE IF NOT EXISTS vip_subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tier_id TEXT NOT NULL REFERENCES vip_tiers(id),
    status TEXT NOT NULL, -- 'active', 'expired', 'cancelled'
    start_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Role Votes (Leaderboard)
CREATE TABLE IF NOT EXISTS role_votes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    period TEXT NOT NULL, -- 'YYYY-MM'
    vote_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(role_id, period)
);

-- Vote Logs (Audit)
CREATE TABLE IF NOT EXISTS vote_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    amount INTEGER NOT NULL,
    period TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Chat Pricing Configuration
CREATE TABLE IF NOT EXISTS chat_pricing (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    model_id UUID REFERENCES models(id) ON DELETE CASCADE, 
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    price_per_msg INTEGER NOT NULL DEFAULT 0, -- Tokens
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(model_id, role_id)
);

-- +goose Down
DROP TABLE chat_pricing;
DROP TABLE vote_logs;
DROP TABLE role_votes;
DROP TABLE vip_subscriptions;
DROP TABLE user_quotas;
DROP TABLE user_assets;
DROP TABLE vip_tiers;
