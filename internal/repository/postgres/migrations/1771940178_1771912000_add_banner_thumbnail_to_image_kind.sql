-- +migrate Up
ALTER TYPE image_kind ADD VALUE IF NOT EXISTS 'banner';
ALTER TYPE image_kind ADD VALUE IF NOT EXISTS 'chapter_thumbnail';

-- +migrate Down
-- Postgres does not support removing enum values; down is a no-op