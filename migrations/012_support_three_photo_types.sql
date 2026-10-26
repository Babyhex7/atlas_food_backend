-- Migration: Support three distinct photo_types (series, range, guide)
-- Created at: 2026-10-26
--
-- Karakteristik tipe foto porsi:
-- 1. 'series': Foto terpisah per porsi (A–H di piring/mangkok terpisah)
-- 2. 'range' : 1 foto memuat 1 jenis makanan yang sama dengan tingkatan/gradasi ukuran (kecil -> sedang -> besar)
-- 3. 'guide' : 1 foto komposit memuat berbagai ragam jenis/varian makanan berbeda (aneka biskuit, kerupuk, dll.)

ALTER TABLE `foods` MODIFY COLUMN `photo_type` ENUM('series', 'range', 'guide') NOT NULL DEFAULT 'series';
