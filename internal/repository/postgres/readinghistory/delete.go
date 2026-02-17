package readinghistory

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) UnmarkAsRead(ctx context.Context, userID uint, chapterID uint) error {
	const op = richerror.Op("repository.postgres.readinghistory.UnmarkAsRead")

	query := `
		DELETE FROM reading_history
		WHERE user_id = $1 AND chapter_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, userID, chapterID)
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to unmark chapter as read")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to get rows affected")
	}

	if rowsAffected == 0 {
		return richerror.New(op).
			WithMessage("reading history not found").
			WithKind(richerror.KindNotFound)
	}

	return nil
}