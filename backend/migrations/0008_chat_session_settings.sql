-- chat session settings & worldbook support for dialog MVP

ALTER TABLE chat_sessions
    ADD COLUMN IF NOT EXISTS mode TEXT NOT NULL DEFAULT 'sfw',
    ADD COLUMN IF NOT EXISTS settings_json JSONB NOT NULL DEFAULT jsonb_build_object(
        'temperature', 0.7,
        'max_tokens', 512,
        'narrative_focus', 'balanced',
        'action_richness', 'medium',
        'sfw_mode', true,
        'immersive', true
    );

CREATE TABLE IF NOT EXISTS worldbooks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    data_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(role_id)
);

