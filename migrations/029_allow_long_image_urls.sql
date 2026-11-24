-- Migration 029: Allow longer URLs for as_served_images to support SVG base64 / data URIs
ALTER TABLE as_served_images MODIFY COLUMN image_url MEDIUMTEXT NOT NULL;
ALTER TABLE as_served_images MODIFY COLUMN thumbnail_url MEDIUMTEXT;
