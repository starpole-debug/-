-- Align chat tables with current backend code (model_key, status, role fields, metadata)

-- chat_sessions adjustments
ALTER TABLE chat_sessions
    ADD COLUMN IF NOT EXISTS model_key TEXT NOT NULL DEFAULT 'mock-fallback',
    ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'active';

-- Best-effort backfill: if legacy model_id stored, mirror it into model_key text
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='chat_sessions' AND column_name='model_id'
    ) THEN
        UPDATE chat_sessions
        SET model_key = COALESCE(model_key, model_id::text, 'mock-fallback');
    END IF;
END $$;

-- chat_messages adjustments
ALTER TABLE chat_messages
    ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'user',
    ADD COLUMN IF NOT EXISTS is_important BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS metadata JSONB NOT NULL DEFAULT '{}'::jsonb;

-- Backfill role from legacy sender_type if present
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='chat_messages' AND column_name='sender_type'
    ) THEN
        UPDATE chat_messages
        SET role = sender_type
        WHERE (sender_type IS NOT NULL AND sender_type <> '');
    END IF;
END $$;

-- Helpful indexes
CREATE INDEX IF NOT EXISTS idx_chat_sessions_updated_at ON chat_sessions(updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_chat_messages_session_id ON chat_messages(session_id);
CREATE INDEX IF NOT EXISTS idx_chat_messages_created_at ON chat_messages(created_at);

-- models table alignment
ALTER TABLE models
    ADD COLUMN IF NOT EXISTS description TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS base_url TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS model_name TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS api_key TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS temperature DOUBLE PRECISION NOT NULL DEFAULT 0.7,
    ADD COLUMN IF NOT EXISTS max_tokens INT NOT NULL DEFAULT 512,
    ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'active',
    ADD COLUMN IF NOT EXISTS max_context_tokens INT NOT NULL DEFAULT 2048;

-- If legacy model_key column exists, copy into model_name for compatibility
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='models' AND column_name='model_key'
    ) THEN
        UPDATE models SET model_name = COALESCE(NULLIF(model_name, ''), model_key);
    END IF;
END $$;
