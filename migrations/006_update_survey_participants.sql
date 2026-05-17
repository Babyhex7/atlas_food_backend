-- Migration: Update survey_participants to require user_id and remove anonymous support
-- Created at: 2026-05-17

-- Hapus data participant anonymous (tanpa user_id) sebelum enforce NOT NULL
DELETE FROM survey_participants
WHERE user_id IS NULL OR user_id = '';

-- Pastikan user_id NOT NULL (guarded)
SET @user_id_nullable := (
    SELECT IS_NULLABLE
    FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'survey_participants'
      AND COLUMN_NAME = 'user_id'
    LIMIT 1
);

SET @sql_user := IF(
    @user_id_nullable = 'YES',
    'ALTER TABLE survey_participants MODIFY COLUMN user_id CHAR(36) NOT NULL',
    'SELECT 1'
);

PREPARE stmt_user FROM @sql_user;
EXECUTE stmt_user;
DEALLOCATE PREPARE stmt_user;

-- Drop kolom is_anonymous jika masih ada (guarded)
SET @col_exists := (
    SELECT COUNT(*)
    FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'survey_participants'
      AND COLUMN_NAME = 'is_anonymous'
);

SET @sql_drop := IF(
    @col_exists > 0,
    'ALTER TABLE survey_participants DROP COLUMN is_anonymous',
    'SELECT 1'
);

PREPARE stmt_drop FROM @sql_drop;
EXECUTE stmt_drop;
DEALLOCATE PREPARE stmt_drop;
