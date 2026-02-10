-- +migrate Up

-- Create series table
CREATE TABLE IF NOT EXISTS series (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug TEXT NOT NULL,          
    slug_id CHAR(8) NOT NULL,     
    full_slug TEXT NOT NULL,      
    description TEXT,
    author VARCHAR(255),
    artist VARCHAR(255),
    status VARCHAR(50) NOT NULL DEFAULT 'ongoing',
    type VARCHAR(50) NOT NULL,
    genres JSONB DEFAULT '[]'::JSONB,
    alternative_titles JSONB DEFAULT '[]'::JSONB,
    cover_image_url TEXT,
    publication_year INTEGER,
    view_count INTEGER NOT NULL DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0.00,
    is_published BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- indexes / constraints (unchanged)
CREATE INDEX idx_series_slug ON series(slug);
CREATE INDEX idx_series_status ON series(status);
CREATE INDEX idx_series_type ON series(type);
CREATE INDEX idx_series_is_published ON series(is_published);
CREATE INDEX idx_series_created_at ON series(created_at DESC);
CREATE INDEX idx_series_rating ON series(rating DESC);
CREATE INDEX idx_series_view_count ON series(view_count DESC);
CREATE INDEX idx_series_genres ON series USING GIN(genres);
CREATE INDEX idx_series_alternative_titles ON series USING GIN(alternative_titles);

ALTER TABLE series ADD CONSTRAINT chk_series_status
    CHECK (status IN ('ongoing', 'completed', 'hiatus', 'cancelled'));

ALTER TABLE series ADD CONSTRAINT chk_series_type
    CHECK (type IN ('manga', 'manhwa', 'manhua', 'comic', 'webtoon'));

ALTER TABLE series ADD CONSTRAINT chk_series_rating
    CHECK (rating >= 0.00 AND rating <= 10.00);

ALTER TABLE series ADD CONSTRAINT chk_series_view_count
    CHECK (view_count >= 0);

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION update_series_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_series_updated_at
    BEFORE UPDATE ON series
    FOR EACH ROW
    EXECUTE FUNCTION update_series_updated_at();
-- +migrate StatementEnd

-- +migrate Down

-- +migrate StatementBegin
DROP TRIGGER IF EXISTS trg_series_updated_at ON series;
DROP FUNCTION IF EXISTS update_series_updated_at();
-- +migrate StatementEnd

DROP TABLE IF EXISTS series;
