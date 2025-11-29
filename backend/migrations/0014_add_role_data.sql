-- Add data column to roles table for builder state
ALTER TABLE roles ADD COLUMN IF NOT EXISTS data JSONB;
