-- +goose Up
-- Ensure a demo role exists for testing / search
INSERT INTO roles (id, creator_id, name, description, avatar_url, tags, abilities, allow_clone, status, role_version)
SELECT
    '11111111-1111-1111-1111-111111111111',
    id,
    'Nebula Demo',
    'A friendly AI guide who helps you explore the platform and start conversations.',
    'https://api.dicebear.com/7.x/bottts/svg?seed=NebulaDemo',
    ARRAY['demo','guide','friendly'],
    ARRAY['onboarding','answering FAQs','giving tips'],
    FALSE,
    'published',
    'v1'
FROM users
WHERE is_admin = true
LIMIT 1
ON CONFLICT (id) DO NOTHING;

-- +goose Down
DELETE FROM roles WHERE id = '11111111-1111-1111-1111-111111111111';
