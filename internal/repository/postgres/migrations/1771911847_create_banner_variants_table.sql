-- +migrate Up
CREATE TABLE banner_variants (
    id          BIGSERIAL PRIMARY KEY,
    series_id   BIGINT NOT NULL REFERENCES series(id) ON DELETE CASCADE,
    kind        VARCHAR(50) NOT NULL,
    image_url   TEXT NOT NULL,
    remote_path TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_banner_variants_series_id ON banner_variants(series_id);

-- +migrate Down
DROP INDEX IF EXISTS idx_banner_variants_series_id;
DROP TABLE IF EXISTS banner_variants;