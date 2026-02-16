-- +migrate Up
-- Create bookmarks table
CREATE TABLE IF NOT EXISTS bookmarks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    series_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign keys
    CONSTRAINT fk_bookmarks_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_bookmarks_series FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE,
    
    -- Unique constraint: one bookmark per user per series
    CONSTRAINT uq_bookmarks_user_series UNIQUE (user_id, series_id)
);

-- Create indexes for fast lookups
CREATE INDEX idx_bookmarks_user_id ON bookmarks(user_id);
CREATE INDEX idx_bookmarks_series_id ON bookmarks(series_id);
CREATE INDEX idx_bookmarks_created_at ON bookmarks(created_at DESC);

-- +migrate Down
DROP TABLE IF EXISTS bookmarks;