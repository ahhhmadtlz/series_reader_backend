package bannervariant

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *BannerVariantRepo) CreateBannerVariant(ctx context.Context, variant entity.BannerVariant) (entity.BannerVariant, error) {
	const op = richerror.Op("bannervariant.CreateBannerVariant")

	query := `
		INSERT INTO banner_variants (series_id, kind, image_url, remote_path, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, series_id, kind, image_url, remote_path, created_at
	`

	var result entity.BannerVariant
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
		return entity.BannerVariant{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to create banner variant").
			WithKind(richerror.KindUnexpected)
	}

	return result, nil
}