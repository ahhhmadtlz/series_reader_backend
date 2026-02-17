package readinghistory

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) IsChapterRead(ctx context.Context, userID uint, chapterID uint) (bool, error) {
	const op = richerror.Op("repository.postgres.readinghistory.IsChapterRead")

	query := `
		SELECT EXISTS(
			SELECT 1 FROM reading_history
			WHERE user_id = $1 AND chapter_id = $2
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID, chapterID).Scan(&exists)

	if err != nil {
		return false, richerror.New(op).
			WithErr(err).
			WithMessage("failed to check if chapter is read")
	}

	return exists, nil
}