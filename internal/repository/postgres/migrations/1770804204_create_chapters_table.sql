-- +migrate Up

CREATE TABLE IF NOT EXISTS chapters (
  id SERIAL PRIMARY KEY,
  series_id INTEGER NOT NULL REFERENCES series(id) ON DELETE CASCADE,
  chapter_number DECIMAL(10,2) NOT NULL,
  title VARCHAR(255),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (series_id, chapter_number)
);

CREATE TABLE IF NOT EXISTS chapter_pages (
  id SERIAL PRIMARY KEY,
  chapter_id INTEGER NOT NULL REFERENCES chapters(id) ON DELETE CASCADE,
  page_number INTEGER NOT NULL,
  image_url TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (chapter_id, page_number)
);

CREATE INDEX IF NOT EXISTS idx_chapters_series_id
  ON chapters(series_id);

CREATE INDEX IF NOT EXISTS idx_chapter_pages_chapter_id
  ON chapter_pages(chapter_id);

CREATE INDEX IF NOT EXISTS idx_chapter_pages_chapter_id_page_number
  ON chapter_pages(chapter_id, page_number);


-- +migrate Down

DROP INDEX IF EXISTS idx_chapter_pages_chapter_id_page_number;
DROP INDEX IF EXISTS idx_chapter_pages_chapter_id;
DROP INDEX IF EXISTS idx_chapters_series_id;

DROP TABLE IF EXISTS chapter_pages;
DROP TABLE IF EXISTS chapters;
