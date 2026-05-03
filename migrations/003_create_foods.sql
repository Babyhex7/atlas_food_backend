-- Migration: Create food database tables
-- Created at: 2024-01-01

-- Tabel categories
CREATE TABLE IF NOT EXISTS categories (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    icon VARCHAR(50),
    display_order INT DEFAULT 0,
    INDEX idx_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabel foods
CREATE TABLE IF NOT EXISTS foods (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    local_name VARCHAR(255),
    description TEXT,
    category_id CHAR(36),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id),
    INDEX idx_name (name),
    FULLTEXT INDEX ft_name (name, local_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabel food_categories (many-to-many)
CREATE TABLE IF NOT EXISTS food_categories (
    food_id CHAR(36) NOT NULL,
    category_id CHAR(36) NOT NULL,
    PRIMARY KEY (food_id, category_id),
    FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabel nutrient_units
CREATE TABLE IF NOT EXISTS nutrient_units (
    id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL,
    symbol VARCHAR(10) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert default units
INSERT INTO nutrient_units (code, name, symbol) VALUES
('kcal', 'Kilokalori', 'kcal'),
('g', 'Gram', 'g'),
('mg', 'Miligram', 'mg'),
('mcg', 'Mikrogram', 'μg')
ON DUPLICATE KEY UPDATE name=name;

-- Tabel nutrient_types
CREATE TABLE IF NOT EXISTS nutrient_types (
    id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(30) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    unit_id INT NOT NULL,
    display_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (unit_id) REFERENCES nutrient_units(id),
    INDEX idx_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert default nutrient types
INSERT INTO nutrient_types (code, name, unit_id, display_order) VALUES
('energy', 'Energi', 1, 1),
('protein', 'Protein', 2, 2),
('carbs', 'Karbohidrat', 2, 3),
('fat', 'Lemak Total', 2, 4),
('fiber', 'Serat', 2, 5),
('sugar', 'Gula', 2, 6),
('sodium', 'Natrium', 3, 7),
('calcium', 'Kalsium', 3, 8),
('iron', 'Zat Besi', 3, 9),
('vit_c', 'Vitamin C', 3, 10)
ON DUPLICATE KEY UPDATE name=name;

-- Tabel food_nutrients
CREATE TABLE IF NOT EXISTS food_nutrients (
    food_id CHAR(36) NOT NULL,
    nutrient_type_id INT NOT NULL,
    value_per_100g DECIMAL(10, 4) NOT NULL,
    PRIMARY KEY (food_id, nutrient_type_id),
    FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE CASCADE,
    FOREIGN KEY (nutrient_type_id) REFERENCES nutrient_types(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
