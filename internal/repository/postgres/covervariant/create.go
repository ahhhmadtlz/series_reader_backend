package covervariant

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *CoverVariantRepo) CreateCoverVariant(ctx context.Context, variant entity.CoverVariant) (entity.CoverVariant, error) {
	const op = richerror.Op("covervariant.CreateCoverVariant")

	query := `
		INSERT INTO cover_variants (series_id, kind, image_url, remote_path, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, series_id, kind, image_url, remote_path, created_at
	`

	var result entity.CoverVariant
	err := r.db.QueryRowContext(ctx, query,
		variant.SeriesID,
		variant.Kind,
		variant.ImageURL,
		variant.RemotePath,
	).Scan(
		&result.ID,
		&result.SeriesID,
		&result.Kind,
		&result.ImageURL,
		&result.RemotePath,
		&result.CreatedAt,
	)
	if err != nil {
		return entity.CoverVariant{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to create cover variant").
			WithKind(richerror.KindUnexpected)
	}

	return result, nil
}