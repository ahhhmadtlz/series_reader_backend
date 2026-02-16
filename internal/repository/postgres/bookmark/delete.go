package bookmark

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) Delete(ctx context.Context, userID uint, seriesID uint) error {
	const op = richerror.Op("repository.postgres.bookmark.Delete")

	query := `DELETE FROM bookmarks WHERE user_id = $1 AND series_id = $2`

	result, err := r.db.ExecContext(ctx, query, userID, seriesID)
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to delete bookmark")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to get rows affected")
	}

	if rowsAffected == 0 {
		return richerror.New(op).
			WithMessage("bookmark not found").
			WithKind(richerror.KindNotFound)
	}

	return nil
}