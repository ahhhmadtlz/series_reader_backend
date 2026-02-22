-- +migrate Up

-- Drop the uploaded_image_id column (old approach we don't use)
ALTER TABLE chapter_pages
    DROP COLUMN IF EXISTS uploaded_image_id;

-- Add image_url back (our approach stores URL directly)
ALTER TABLE chapter_pages
    ADD COLUMN IF NOT EXISTS image_url TEXT NOT NULL DEFAULT '';

-- +migrate Down

ALTER TABLE chapter_pages
    DROP COLUMN IF EXISTS image_url;

ALTER TABLE chapter_pages
    ADD COLUMN IF NOT EXISTS uploaded_image_id INTEGER
    REFERENCES uploaded_images(id) ON DELETE CASCADE;