package bannervariant

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *BannerVariantRepo) GetBannerVariantsBySeriesID(ctx context.Context, seriesID uint) ([]entity.BannerVariant, error) {
	const op = richerror.Op("bannervariant.GetBannerVariantsBySeriesID")

	query := `
		SELECT id, series_id, kind, image_url, remote_path, created_at
		FROM banner_variants
		WHERE series_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, seriesID)
	if err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get banner variants").
			WithKind(richerror.KindUnexpected)
	}
	defer rows.Close()

	var variants []entity.BannerVariant
	for rows.Next() {
		var v entity.BannerVariant
		if err := rows.Scan(&v.ID, &v.SeriesID, &v.Kind, &v.ImageURL, &v.RemotePath, &v.CreatedAt); err != nil {
			return nil, richerror.New(op).
				WithErr(err).
				WithMessage("failed to scan banner variant").
				WithKind(richerror.KindUnexpected)
		}
		variants = append(variants, v)
	}

	return variants, nil
}