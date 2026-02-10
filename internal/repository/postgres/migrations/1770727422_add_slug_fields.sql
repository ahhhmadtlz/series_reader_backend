-- +migrate Up
ALTER TABLE series ADD COLUMN IF NOT EXISTS slug_id VARCHAR(50);
ALTER TABLE series ADD COLUMN IF NOT EXISTS full_slug VARCHAR(300);

UPDATE series SET slug_id = slug, full_slug = slug WHERE slug_id IS NULL;

-- +migrate Down
ALTER TABLE series DROP COLUMN IF EXISTS full_slug;
ALTER TABLE series DROP COLUMN IF EXISTS slug_id;