package imagevariant

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) DeleteVariantsByPageID(ctx context.Context, pageID uint) error {
	const op = richerror.Op("postgres.imagevariant.DeleteVariantsByPageID")

	const query = `DELETE FROM image_variants WHERE chapter_page_id = $1`

	if _, err := r.db.ExecContext(ctx, query, pageID); err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to delete image variants").
			WithKind(richerror.KindUnexpected)
	}

	return nil
}