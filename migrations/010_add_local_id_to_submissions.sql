-- Migration: Add local_id to survey_submissions for offline idempotency
-- Created at: 2026-09-02
-- Purpose: Agar backend bisa mendeteksi submission duplikat dari offline sync (Idempotency Key)

ALTER TABLE survey_submissions
  ADD COLUMN local_id VARCHAR(36) NULL UNIQUE COMMENT 'UUID dari client (offline localId), untuk idempotency duplicate detection',
  ADD COLUMN total_energy DECIMAL(10,2) DEFAULT 0 COMMENT 'Total energi harian (kcal)',
  ADD COLUMN total_protein DECIMAL(10,2) DEFAULT 0 COMMENT 'Total protein harian (gram)',
  ADD COLUMN total_carbs DECIMAL(10,2) DEFAULT 0 COMMENT 'Total karbohidrat harian (gram)',
  ADD COLUMN total_fat DECIMAL(10,2) DEFAULT 0 COMMENT 'Total lemak harian (gram)';

-- Index untuk lookup cepat saat idempotency check
CREATE INDEX idx_local_id ON survey_submissions (local_id);
