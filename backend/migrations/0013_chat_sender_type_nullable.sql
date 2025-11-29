-- align sender_type with role and make it nullable to avoid insert failures
ALTER TABLE chat_messages
    ADD COLUMN IF NOT EXISTS sender_type TEXT;

-- backfill any null sender_type using role
UPDATE chat_messages
SET sender_type = COALESCE(sender_type, role, 'user')
WHERE sender_type IS NULL;

-- relax constraint so new inserts without sender_type don't fail
ALTER TABLE chat_messages
    ALTER COLUMN sender_type DROP NOT NULL,
    ALTER COLUMN sender_type DROP DEFAULT;
