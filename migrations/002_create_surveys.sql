-- Migration: Create surveys and related tables
-- Created at: 2024-01-01

-- Tabel locales untuk multi-language (hardcoded: id, en)
CREATE TABLE IF NOT EXISTS locales (
    id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert default locales
INSERT INTO locales (code, name) VALUES
('id', 'Indonesia'),
('en', 'English')
ON DUPLICATE KEY UPDATE name=name;

-- Tabel surveys
CREATE TABLE IF NOT EXISTS surveys (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    slug VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    meals_config JSON NOT NULL,
    prompts JSON,
    locale_id INT DEFAULT 1,
    start_date DATE,
    end_date DATE,
    status ENUM('draft', 'active', 'closed') DEFAULT 'draft',
    access_token VARCHAR(255),
    created_by CHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (locale_id) REFERENCES locales(id),
    INDEX idx_slug (slug),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabel survey_participants untuk respondent terdaftar
CREATE TABLE IF NOT EXISTS survey_participants (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    survey_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    alias VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (survey_id) REFERENCES surveys(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_survey (survey_id),
    INDEX idx_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
