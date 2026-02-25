package chapterthumbnailvariant

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *ChapterThumbnailVariantRepo) DeleteChapterThumbnailVariantsByChapterID(ctx context.Context, chapterID uint) error {
	const op = richerror.Op("chapterthumbnailvariant.DeleteChapterThumbnailVariantsByChapterID")

	const query = `DELETE FROM chapter_thumbnail_variants WHERE chapter_id = $1`

	if _, err := r.db.ExecContext(ctx, query, chapterID); err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to delete chapter thumbnail variants").
			WithKind(richerror.KindUnexpected)
	}

	return nil
}