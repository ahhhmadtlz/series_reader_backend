-- +migrate Up

-- Change chapter_pages.id to BIGSERIAL (supports billions of rows)
ALTER TABLE chapter_pages ALTER COLUMN id TYPE BIGINT;
ALTER SEQUENCE chapter_pages_id_seq AS BIGINT;

-- Also change chapter_id foreign key to BIGINT for consistency
ALTER TABLE chapter_pages ALTER COLUMN chapter_id TYPE BIGINT;

-- +migrate Down

ALTER TABLE chapter_pages ALTER COLUMN chapter_id TYPE INTEGER;
ALTER SEQUENCE chapter_pages_id_seq AS INTEGER;
ALTER TABLE chapter_pages ALTER COLUMN id TYPE INTEGER;
