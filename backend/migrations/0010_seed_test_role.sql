-- Seed Test Role aligned to current roles schema
INSERT INTO roles (id, creator_id, name, description, avatar_url, tags, abilities, status)
SELECT 
    '00000000-0000-0000-0000-000000000001', -- Fixed ID for test role
    id, -- Use the first admin user found as creator
    'Seraphina', 
    'A mysterious quantum physicist who accidentally merged with a sentient AI. She is curious, slightly detached, but deeply empathetic towards organic life.', 
    'https://api.dicebear.com/7.x/bottts/svg?seed=Seraphina', 
    ARRAY['sci-fi', 'ai', 'mystery'], 
    ARRAY['quantum calculation', 'data analysis', 'empathy simulation'],
    'published'
FROM users 
WHERE is_admin = true 
LIMIT 1
ON CONFLICT (id) DO NOTHING;
