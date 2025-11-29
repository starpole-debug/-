-- Add per-call coin price for models
ALTER TABLE models ADD COLUMN IF NOT EXISTS price_coins BIGINT NOT NULL DEFAULT 0;
UPDATE models SET price_coins = 0 WHERE price_coins IS NULL;
