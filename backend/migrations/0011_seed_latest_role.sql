-- Ensure we have a published role for testing "Latest"
INSERT INTO roles (id, creator_id, name, description, avatar_url, tags, abilities, status, updated_at)
SELECT 
    '00000000-0000-0000-0000-000000000002', -- Different ID to avoid conflict with previous test if needed, or use same to update
    id, 
    'Nova', 
    'A futuristic cyber-enhanced explorer from the year 3024. She loves neon lights and old-school jazz.', 
    'https://api.dicebear.com/7.x/bottts/svg?seed=Nova', 
    ARRAY['cyberpunk', 'music', 'explorer'], 
    ARRAY['hacking', 'navigation', 'history'],
    'published',
    NOW()
FROM users 
WHERE is_admin = true 
LIMIT 1
ON CONFLICT (id) DO UPDATE SET
    status = 'published',
    updated_at = NOW();
