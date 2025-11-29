-- Image generation providers (NovelAI etc.)
CREATE TABLE IF NOT EXISTS image_providers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    base_url TEXT NOT NULL,
    api_key TEXT NOT NULL,
    max_concurrency INT NOT NULL DEFAULT 5,
    weight INT NOT NULL DEFAULT 1,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_image_providers_status ON image_providers(status);

-- Prompt presets (JSON) for rendering instructions
CREATE TABLE IF NOT EXISTS image_presets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    preset_json JSONB NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_image_presets_status ON image_presets(status);

-- Jobs for user generated images from chat
CREATE TABLE IF NOT EXISTS image_jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id UUID NOT NULL REFERENCES chat_sessions(id) ON DELETE CASCADE,
    message_id UUID,
    provider_id UUID NOT NULL REFERENCES image_providers(id),
    preset_id UUID,
    prompt TEXT NOT NULL,
    negative_prompt TEXT,
    final_prompt TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    result_url TEXT,
    error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_image_jobs_user ON image_jobs(user_id);
CREATE INDEX IF NOT EXISTS idx_image_jobs_session ON image_jobs(session_id);
CREATE INDEX IF NOT EXISTS idx_image_jobs_status ON image_jobs(status);
