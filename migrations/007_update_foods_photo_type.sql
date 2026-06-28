-- Menambahkan kolom photo_type pada tabel foods
ALTER TABLE `foods` ADD COLUMN `photo_type` ENUM('series', 'range') DEFAULT 'series' AFTER `description`;

-- Menambahkan FULLTEXT INDEX untuk pencarian berkinerja tinggi
ALTER TABLE `foods` ADD FULLTEXT INDEX `idx_food_search` (`name`, `local_name`);
