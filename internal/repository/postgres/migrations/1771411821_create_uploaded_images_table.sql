-- +migrate Up

-- Create image_kind enum
CREATE TYPE image_kind AS ENUM ('avatar', 'cover', 'chapter_page');

-- Create uploaded_images table
CREATE TABLE IF NOT EXISTS uploaded_images (
    id           SERIAL PRIMARY KEY,
    owner_id     INTEGER NOT NULL,
    kind         image_kind NOT NULL,
    filename     VARCHAR(255) NOT NULL,
    stored_path  VARCHAR(500) NOT NULL UNIQUE,
    url          VARCHAR(500) NOT NULL,
    mime_type    VARCHAR(100) NOT NULL,
    size_bytes   BIGINT NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_size_positive CHECK (size_bytes >= 0),
    CONSTRAINT chk_valid_mime CHECK (
        mime_type IN (
            'image/jpeg',
            'image/jpg',
            'image/png',
            'image/webp',
            'image/gif'
        )
    )
);

-- Indexes for common queries
CREATE INDEX idx_uploaded_images_owner_kind ON uploaded_images(owner_id, kind);
CREATE INDEX idx_uploaded_images_kind ON uploaded_images(kind);
CREATE INDEX idx_uploaded_images_stored_path ON uploaded_images(stored_path);
CREATE INDEX idx_uploaded_images_created_at ON uploaded_images(created_at DESC);

-- Enforce single active avatar/cover per owner
CREATE UNIQUE INDEX uq_uploaded_images_single_avatar 
    ON uploaded_images(owner_id) 
    WHERE kind = 'avatar';

CREATE UNIQUE INDEX uq_uploaded_images_single_cover 
    ON uploaded_images(owner_id) 
    WHERE kind = 'cover';

-- Step 1: Add new column
ALTER TABLE chapter_pages 
    ADD COLUMN uploaded_image_id INTEGER 
    REFERENCES uploaded_images(id) ON DELETE CASCADE;

-- Step 2: Migrate existing data
-- +migrate StatementBegin
DO $MIGRATION$
BEGIN
    IF EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'chapter_pages'
          AND column_name = 'image_url'
    ) THEN

        INSERT INTO uploaded_images (
            owner_id,
            kind,
            filename,
            stored_path,
            url,
            mime_type,
            size_bytes
        )
        SELECT 
            cp.chapter_id,
            'chapter_page'::image_kind,
            'migrated_page_' || cp.id || '.jpg',
            cp.image_url,
            cp.image_url,
            'image/jpeg',
            0
        FROM chapter_pages cp
        WHERE cp.image_url IS NOT NULL 
          AND cp.image_url != ''
        ON CONFLICT (stored_path) DO NOTHING;

        UPDATE chapter_pages cp
        SET uploaded_image_id = ui.id
        FROM uploaded_images ui
        WHERE ui.stored_path = cp.image_url
          AND ui.kind = 'chapter_page'
          AND cp.image_url IS NOT NULL;

    END IF;
END;
$MIGRATION$;
-- +migrate StatementEnd

-- Step 3: Drop old column
ALTER TABLE chapter_pages 
    DROP COLUMN IF EXISTS image_url;

-- Step 4: Index for lookups
CREATE INDEX idx_chapter_pages_uploaded_image 
    ON chapter_pages(uploaded_image_id);


-- +migrate Down

ALTER TABLE chapter_pages 
    ADD COLUMN image_url VARCHAR(500);

UPDATE chapter_pages cp
SET image_url = ui.url
FROM uploaded_images ui
WHERE cp.uploaded_image_id = ui.id;

DROP INDEX IF EXISTS idx_chapter_pages_uploaded_image;
ALTER TABLE chapter_pages DROP COLUMN IF EXISTS uploaded_image_id;

DELETE FROM uploaded_images WHERE kind = 'chapter_page';

DROP INDEX IF EXISTS uq_uploaded_images_single_cover;
DROP INDEX IF EXISTS uq_uploaded_images_single_avatar;
DROP INDEX IF EXISTS idx_uploaded_images_created_at;
DROP INDEX IF EXISTS idx_uploaded_images_stored_path;
DROP INDEX IF EXISTS idx_uploaded_images_kind;
DROP INDEX IF EXISTS idx_uploaded_images_owner_kind;

DROP TABLE IF EXISTS uploaded_images;

DROP TYPE IF EXISTS image_kind;