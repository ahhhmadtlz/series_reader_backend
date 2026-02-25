package imagevariant

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) GetVariantsByPageID(ctx context.Context, pageID uint) ([]entity.ImageVariant, error) {
	const op = richerror.Op("postgres.imagevariant.GetVariantsByPageID")

	const query = `
		SELECT id, chapter_page_id, kind, image_url, remote_path, created_at
		FROM image_variants
		WHERE chapter_page_id = $1
		ORDER BY kind`

	rows, err := r.db.QueryContext(ctx, query, pageID)
	if err != nil {
		logger.Error("failed to query image variants", "error", err)
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get image variants").
			WithKind(richerror.KindUnexpected)
	}
	defer rows.Close()

	var variants []entity.ImageVariant
	for rows.Next() {
		var v entity.ImageVariant
		if err := rows.Scan(&v.ID, &v.ChapterPageID, &v.Kind, &v.ImageURL, &v.RemotePath, &v.CreatedAt); err != nil {
			return nil, richerror.New(op).
				WithErr(err).
				WithMessage("failed to scan image variant").
				WithKind(richerror.KindUnexpected)
		}
		variants = append(variants, v)
	}

	if err := rows.Err(); err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("rows error while reading image variants").
			WithKind(richerror.KindUnexpected)
	}

	return variants, nil
}