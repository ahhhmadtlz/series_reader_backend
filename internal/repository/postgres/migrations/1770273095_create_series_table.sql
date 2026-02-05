-- +migrate Up
-- Create series table
CREATE TABLE IF NOT EXISTS series (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
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

-- Add indexes for better query performance
CREATE INDEX idx_series_slug ON series(slug);
CREATE INDEX idx_series_status ON series(status);
CREATE INDEX idx_series_type ON series(type);
CREATE INDEX idx_series_is_published ON series(is_published);
CREATE INDEX idx_series_created_at ON series(created_at DESC);
CREATE INDEX idx_series_rating ON series(rating DESC);
CREATE INDEX idx_series_view_count ON series(view_count DESC);

-- GIN indexes for JSONB columns to enable efficient searching
CREATE INDEX idx_series_genres ON series USING GIN(genres);
CREATE INDEX idx_series_alternative_titles ON series USING GIN(alternative_titles);

-- Add constraint to validate status values
ALTER TABLE series ADD CONSTRAINT chk_series_status 
    CHECK (status IN ('ongoing', 'completed', 'hiatus', 'cancelled'));

-- Add constraint to validate type values
ALTER TABLE series ADD CONSTRAINT chk_series_type 
    CHECK (type IN ('manga', 'manhwa', 'manhua', 'comic', 'webtoon'));

-- Add constraint for rating range (0.00 to 10.00)
ALTER TABLE series ADD CONSTRAINT chk_series_rating 
    CHECK (rating >= 0.00 AND rating <= 10.00);

-- Add constraint for view_count (must be non-negative)
ALTER TABLE series ADD CONSTRAINT chk_series_view_count 
    CHECK (view_count >= 0);

-- Function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_series_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to call the function before any UPDATE on series table
CREATE TRIGGER trg_series_updated_at
    BEFORE UPDATE ON series
    FOR EACH ROW
    EXECUTE FUNCTION update_series_updated_at();

-- +migrate Down
-- Drop trigger and function
DROP TRIGGER IF EXISTS trg_series_updated_at ON series;
DROP FUNCTION IF EXISTS update_series_updated_at();

-- Drop table (this will also drop all indexes and constraints)
DROP TABLE IF EXISTS series;