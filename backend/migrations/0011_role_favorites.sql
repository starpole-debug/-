-- user-role favorites for quick access in chat
CREATE TABLE IF NOT EXISTS role_favorites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, role_id)
);

CREATE INDEX IF NOT EXISTS idx_role_favorites_user ON role_favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_role_favorites_role ON role_favorites(role_id);
