-- Migration: Create submissions table
-- Created at: 2024-01-01

-- Tabel survey_submissions
CREATE TABLE IF NOT EXISTS survey_submissions (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    survey_id CHAR(36) NOT NULL,
    participant_id CHAR(36),
    respondent_name VARCHAR(255),
    respondent_email VARCHAR(255),
    meals_data JSON NOT NULL,
    missing_foods JSON,
    submitted_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (survey_id) REFERENCES surveys(id) ON DELETE CASCADE,
    FOREIGN KEY (participant_id) REFERENCES survey_participants(id) ON DELETE SET NULL,
    INDEX idx_survey (survey_id),
    INDEX idx_participant (participant_id),
    INDEX idx_submitted_at (submitted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
