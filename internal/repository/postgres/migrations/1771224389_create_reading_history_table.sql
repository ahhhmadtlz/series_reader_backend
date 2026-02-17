-- +migrate Up
-- Create reading_history table
CREATE TABLE IF NOT EXISTS reading_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    chapter_id INTEGER NOT NULL,
    read_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign keys
    CONSTRAINT fk_reading_history_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_reading_history_chapter FOREIGN KEY (chapter_id) REFERENCES chapters(id) ON DELETE CASCADE,
    
    -- Unique constraint: one reading record per user per chapter
    CONSTRAINT uq_reading_history_user_chapter UNIQUE (user_id, chapter_id)
);

-- Create indexes for fast lookups
CREATE INDEX idx_reading_history_user_id ON reading_history(user_id);
CREATE INDEX idx_reading_history_chapter_id ON reading_history(chapter_id);
CREATE INDEX idx_reading_history_read_at ON reading_history(read_at DESC);
CREATE INDEX idx_reading_history_user_read_at ON reading_history(user_id, read_at DESC);

-- +migrate Down
DROP TABLE IF EXISTS reading_history;