package chapterthumbnailvariant

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *ChapterThumbnailVariantRepo) GetChapterThumbnailVariantsByChapterID(ctx context.Context, chapterID uint) ([]entity.ChapterThumbnailVariant, error) {
	const op = richerror.Op("chapterthumbnailvariant.GetChapterThumbnailVariantsByChapterID")

	query := `
		SELECT id, chapter_id, kind, image_url, remote_path, created_at
		FROM chapter_thumbnail_variants
		WHERE chapter_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, chapterID)
	if err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get chapter thumbnail variants").
			WithKind(richerror.KindUnexpected)
	}
	defer rows.Close()

	var variants []entity.ChapterThumbnailVariant
	for rows.Next() {
		var v entity.ChapterThumbnailVariant
		if err := rows.Scan(&v.ID, &v.ChapterID, &v.Kind, &v.ImageURL, &v.RemotePath, &v.CreatedAt); err != nil {
			return nil, richerror.New(op).
				WithErr(err).
				WithMessage("failed to scan chapter thumbnail variant").
				WithKind(richerror.KindUnexpected)
		}
		variants = append(variants, v)
	}

	return variants, nil
}