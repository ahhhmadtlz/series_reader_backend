package bannervariant

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *BannerVariantRepo) DeleteBannerVariantsBySeriesID(ctx context.Context, seriesID uint) error {
	const op = richerror.Op("bannervariant.DeleteBannerVariantsBySeriesID")

	const query = `DELETE FROM banner_variants WHERE series_id = $1`

	if _, err := r.db.ExecContext(ctx, query, seriesID); err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to delete banner variants").
			WithKind(richerror.KindUnexpected)
	}

	return nil
}