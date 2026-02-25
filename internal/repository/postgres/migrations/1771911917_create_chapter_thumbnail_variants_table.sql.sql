-- +migrate Up
CREATE TABLE chapter_thumbnail_variants (
    id          BIGSERIAL PRIMARY KEY,
    chapter_id  BIGINT NOT NULL REFERENCES chapters(id) ON DELETE CASCADE,
    kind        VARCHAR(50) NOT NULL,
    image_url   TEXT NOT NULL,
    remote_path TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_chapter_thumbnail_variants_chapter_id ON chapter_thumbnail_variants(chapter_id);

-- +migrate Down
DROP INDEX IF EXISTS idx_chapter_thumbnail_variants_chapter_id;
DROP TABLE IF EXISTS chapter_thumbnail_variants;