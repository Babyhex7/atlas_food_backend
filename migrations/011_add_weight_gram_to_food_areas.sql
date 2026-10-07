-- Migration: Add weight_gram to food_areas
-- Created at: 2026-09-30

ALTER TABLE food_areas
ADD COLUMN weight_gram DECIMAL(10, 2) NULL AFTER z_index;
