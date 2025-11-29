ALTER TABLE image_providers
ADD COLUMN IF NOT EXISTS selected_model TEXT DEFAULT '';
