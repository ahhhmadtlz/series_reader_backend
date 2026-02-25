package covervariant

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *CoverVariantRepo) DeleteCoverVariantsBySeriesID(ctx context.Context, seriesID uint) error {
	const op = richerror.Op("covervariant.DeleteCoverVariantsBySeriesID")

	const query = `DELETE FROM cover_variants WHERE series_id = $1`

	if _, err := r.db.ExecContext(ctx, query, seriesID); err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to delete cover variants").
			WithKind(richerror.KindUnexpected)
	}

	return nil
}