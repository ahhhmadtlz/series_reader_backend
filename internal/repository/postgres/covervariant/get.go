package covervariant

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *CoverVariantRepo) GetCoverVariantsBySeriesID(ctx context.Context, seriesID uint) ([]entity.CoverVariant, error) {
	const op = richerror.Op("covervariant.GetCoverVariantsBySeriesID")

	query := `
		SELECT id, series_id, kind, image_url, remote_path, created_at
		FROM cover_variants
		WHERE series_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, seriesID)
	if err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get cover variants").
			WithKind(richerror.KindUnexpected)
	}
	defer rows.Close()

	var variants []entity.CoverVariant
	for rows.Next() {
		var v entity.CoverVariant
		if err := rows.Scan(&v.ID, &v.SeriesID, &v.Kind, &v.ImageURL, &v.RemotePath, &v.CreatedAt); err != nil {
			return nil, richerror.New(op).
				WithErr(err).
				WithMessage("failed to scan cover variant").
				WithKind(richerror.KindUnexpected)
		}
		variants = append(variants, v)
	}

	return variants, nil
}