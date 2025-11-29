-- Add username column and ensure unique logins
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS username TEXT;

UPDATE users
SET username = COALESCE(NULLIF(username, ''), email)
WHERE username IS NULL OR username = '';

ALTER TABLE users
    ALTER COLUMN username SET NOT NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'users_username_key'
    ) THEN
        ALTER TABLE users
            ADD CONSTRAINT users_username_key UNIQUE (username);
    END IF;
END $$;

-- Rename legacy model_key column when present
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'models' AND column_name = 'model_key'
    ) THEN
        ALTER TABLE models RENAME COLUMN model_key TO model_name;
    END IF;
END $$;

ALTER TABLE models
    ADD COLUMN IF NOT EXISTS description TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS base_url TEXT NOT NULL DEFAULT 'https://api.openai.com/v1',
    ADD COLUMN IF NOT EXISTS api_key TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS temperature DOUBLE PRECISION NOT NULL DEFAULT 1.0,
    ADD COLUMN IF NOT EXISTS max_tokens INT NOT NULL DEFAULT 2048,
    ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'active';

ALTER TABLE models
    ADD COLUMN IF NOT EXISTS model_name TEXT NOT NULL DEFAULT '';

UPDATE models
SET model_name = CASE
        WHEN COALESCE(model_name, '') = '' THEN name
        ELSE model_name
    END;

ALTER TABLE models
    ALTER COLUMN model_name SET NOT NULL;

UPDATE models
SET status = CASE WHEN is_enabled THEN 'active' ELSE 'inactive' END;

CREATE INDEX IF NOT EXISTS idx_models_status ON models(status);
