-- Migration: Add profile fields to users
-- Created at: 2026-07-19

ALTER TABLE users
  ADD COLUMN phone VARCHAR(20) NULL AFTER name,
  ADD COLUMN gender ENUM('male', 'female') NULL AFTER phone,
  ADD COLUMN birth_date DATE NULL AFTER gender,
  ADD COLUMN photo_url VARCHAR(500) NULL AFTER birth_date;
