ALTER TABLE image_presets
ADD COLUMN IF NOT EXISTS prompt_model_key TEXT DEFAULT '';
