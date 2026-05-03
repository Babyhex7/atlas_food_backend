-- Migration: Create portion size tables
-- Created at: 2024-01-01

-- Tabel associated_foods (makanan yang sering dipakai bersama)
CREATE TABLE IF NOT EXISTS associated_foods (
    id INT AUTO_INCREMENT PRIMARY KEY,
    food_id CHAR(36) NOT NULL,
    associated_food_id CHAR(36) NOT NULL,
    priority INT DEFAULT 0,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE CASCADE,
    FOREIGN KEY (associated_food_id) REFERENCES foods(id) ON DELETE CASCADE,
    UNIQUE KEY unique_association (food_id, associated_food_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabel food_portion_size_methods
CREATE TABLE IF NOT EXISTS food_portion_size_methods (
    id INT AUTO_INCREMENT PRIMARY KEY,
    food_id CHAR(36) NOT NULL,
    method_type ENUM('as_served', 'guide_image', 'weight') NOT NULL,
    label VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    image_url VARCHAR(500),
    thumbnail_url VARCHAR(500),
    config JSON,
    display_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE CASCADE,
    INDEX idx_food (food_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabel as_served_sets
CREATE TABLE IF NOT EXISTS as_served_sets (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    food_id CHAR(36),
    category VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE SET NULL,
    INDEX idx_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabel as_served_images
CREATE TABLE IF NOT EXISTS as_served_images (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    set_id CHAR(36) NOT NULL,
    label VARCHAR(50) NOT NULL,
    image_url VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500),
    weight_gram DECIMAL(10, 2) NOT NULL,
    description VARCHAR(255),
    display_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (set_id) REFERENCES as_served_sets(id) ON DELETE CASCADE,
    INDEX idx_set (set_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
