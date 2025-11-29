-- +goose Up
-- Follows (User to User)
CREATE TABLE IF NOT EXISTS follows (
    follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (follower_id, following_id)
);
CREATE INDEX IF NOT EXISTS idx_follows_follower ON follows(follower_id);
CREATE INDEX IF NOT EXISTS idx_follows_following ON follows(following_id);

-- Role Enhancements
ALTER TABLE roles ADD COLUMN IF NOT EXISTS background_url TEXT NOT NULL DEFAULT '';

-- Role Version Enhancements
ALTER TABLE role_versions ADD COLUMN IF NOT EXISTS opening_line TEXT NOT NULL DEFAULT '';
ALTER TABLE role_versions ADD COLUMN IF NOT EXISTS tts_config JSONB NOT NULL DEFAULT '{}'::jsonb;

-- +goose Down
ALTER TABLE role_versions DROP COLUMN IF EXISTS tts_config;
ALTER TABLE role_versions DROP COLUMN IF EXISTS opening_line;
ALTER TABLE roles DROP COLUMN IF EXISTS background_url;
DROP TABLE follows;
