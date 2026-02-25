-- +migrate Up

CREATE TYPE image_variant_kind AS ENUM (
    'webp',
    'thumbnail',
    'optimized',
    'cdn'
);

CREATE TABLE image_variants (
    id              BIGSERIAL PRIMARY KEY,
    chapter_page_id BIGINT NOT NULL REFERENCES chapter_pages(id) ON DELETE CASCADE,
    kind            image_variant_kind NOT NULL,
    image_url       TEXT NOT NULL,
    remote_path     TEXT NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (chapter_page_id, kind)
);

CREATE INDEX idx_image_variants_chapter_page_id ON image_variants(chapter_page_id);

-- +migrate Down

DROP INDEX IF EXISTS idx_image_variants_chapter_page_id;
DROP TABLE IF EXISTS image_variants;
DROP TYPE IF EXISTS image_variant_kind;