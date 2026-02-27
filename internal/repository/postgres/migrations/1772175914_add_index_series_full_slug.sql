-- +migrate Up

CREATE INDEX idx_series_full_slug ON series(full_slug);

-- +migrate Down

DROP INDEX IF EXISTS idx_series_full_slug;