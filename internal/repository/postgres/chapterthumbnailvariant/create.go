package chapterthumbnailvariant

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *ChapterThumbnailVariantRepo) CreateChapterThumbnailVariant(ctx context.Context, variant entity.ChapterThumbnailVariant) (entity.ChapterThumbnailVariant, error) {
	const op = richerror.Op("chapterthumbnailvariant.CreateChapterThumbnailVariant")

	query := `
		INSERT INTO chapter_thumbnail_variants (chapter_id, kind, image_url, remote_path, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, chapter_id, kind, image_url, remote_path, created_at
	`

	var result entity.ChapterThumbnailVariant
	err := r.db.QueryRowContext(ctx, query,
		variant.ChapterID,
		variant.Kind,
		variant.ImageURL,
		variant.RemotePath,
	).Scan(
		&result.ID,
		&result.ChapterID,
		&result.Kind,
		&result.ImageURL,
		&result.RemotePath,
		&result.CreatedAt,
	)
	if err != nil {
		return entity.ChapterThumbnailVariant{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to create chapter thumbnail variant").
			WithKind(richerror.KindUnexpected)
	}

	return result, nil
}