package chapterthumbnailvariant

import "database/sql"

type ChapterThumbnailVariantRepo struct {
	db *sql.DB
}

func New(db *sql.DB) *ChapterThumbnailVariantRepo {
	return &ChapterThumbnailVariantRepo{db: db}
}