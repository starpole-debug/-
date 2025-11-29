-- Align chat_sessions schema with repository code (use `settings` column while keeping legacy data).

-- Ensure the canonical settings column exists.
ALTER TABLE chat_sessions
    ADD COLUMN IF NOT EXISTS settings JSONB NOT NULL DEFAULT jsonb_build_object(
        'temperature', 0.7,
        'max_tokens', 512,
        'narrative_focus', 'balanced',
        'action_richness', 'medium',
        'sfw_mode', true,
        'immersive', true
    );

-- Backfill/keep in sync from the legacy settings_json column when present.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'chat_sessions' AND column_name = 'settings_json'
    ) THEN
        UPDATE chat_sessions
        SET settings = settings_json
        WHERE settings_json IS NOT NULL
          AND (settings IS NULL OR settings::text <> settings_json::text);
    END IF;
END $$;

-- Re-assert constraints in case the column already existed without defaults.
ALTER TABLE chat_sessions
    ALTER COLUMN settings SET NOT NULL,
    ALTER COLUMN settings SET DEFAULT jsonb_build_object(
        'temperature', 0.7,
        'max_tokens', 512,
        'narrative_focus', 'balanced',
        'action_richness', 'medium',
        'sfw_mode', true,
        'immersive', true
    );
