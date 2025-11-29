-- +goose Up
CREATE TABLE IF NOT EXISTS community_post_views (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id UUID NOT NULL REFERENCES community_posts(id) ON DELETE CASCADE,
    view_count INTEGER NOT NULL DEFAULT 1,
    last_viewed TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, post_id)
);
CREATE INDEX IF NOT EXISTS idx_post_views_user_time ON community_post_views(user_id, last_viewed DESC);

-- +goose Down
DROP TABLE IF EXISTS community_post_views;

