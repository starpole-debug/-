-- Create presets table
CREATE TABLE IF NOT EXISTS presets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    model_key TEXT, -- Optional binding to a specific model
    blocks JSONB NOT NULL DEFAULT '[]', -- The prompt blocks
    gen_params JSONB, -- Generation parameters (temp, top_p, etc.)
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_presets_creator ON presets(creator_id);
CREATE INDEX IF NOT EXISTS idx_presets_public ON presets(is_public);
