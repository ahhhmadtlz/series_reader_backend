-- +migrate Up

ALTER TABLE chapter_pages
    ADD COLUMN IF NOT EXISTS remote_path TEXT NOT NULL DEFAULT '';

-- +migrate Down

ALTER TABLE chapter_pages
    DROP COLUMN IF EXISTS remote_path;
