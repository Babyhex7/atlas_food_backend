-- Migration 013: Add AI generator tracking fields
ALTER TABLE foods ADD COLUMN IF NOT EXISTS is_ai_generated BOOLEAN DEFAULT FALSE;
ALTER TABLE as_served_images ADD COLUMN IF NOT EXISTS is_ai_generated BOOLEAN DEFAULT FALSE;
