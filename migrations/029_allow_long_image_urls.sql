-- Migration 029: Allow longer URLs for as_served_images and food_portion_size_methods
ALTER TABLE as_served_images MODIFY COLUMN image_url MEDIUMTEXT NOT NULL;
ALTER TABLE as_served_images MODIFY COLUMN thumbnail_url MEDIUMTEXT;
ALTER TABLE food_portion_size_methods MODIFY COLUMN image_url MEDIUMTEXT;
ALTER TABLE food_portion_size_methods MODIFY COLUMN thumbnail_url MEDIUMTEXT;
