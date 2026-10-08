-- Migration: Extend ai_result_logs untuk hasil akhir, versi prompt, dan rincian token
-- Created at: 2026-10-06
-- Purpose: result_data menyimpan hasil yang dikirim ke klien (narasi LLM + perhitungan server),
--          prompt_version membuat hasil bisa ditelusuri ke prompt pembuatnya,
--          updated_at mencatat kapan analisis ulang terakhir dilakukan.

ALTER TABLE ai_result_logs
  ADD COLUMN result_data JSON NULL AFTER raw_response,
  ADD COLUMN prompt_version VARCHAR(20) NOT NULL DEFAULT 'v1' AFTER model_used,
  ADD COLUMN prompt_tokens INT NULL AFTER token_used,
  ADD COLUMN completion_tokens INT NULL AFTER prompt_tokens,
  ADD COLUMN updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP AFTER created_at,
  MODIFY COLUMN model_used VARCHAR(100) NOT NULL DEFAULT '';

UPDATE ai_result_logs SET updated_at = created_at WHERE updated_at IS NULL;
