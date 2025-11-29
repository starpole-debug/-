-- Add per-call share percentages for model revenue split
ALTER TABLE models ADD COLUMN IF NOT EXISTS share_role_pct DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE models ADD COLUMN IF NOT EXISTS share_preset_pct DOUBLE PRECISION NOT NULL DEFAULT 0;
