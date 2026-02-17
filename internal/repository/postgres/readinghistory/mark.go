package readinghistory

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/readinghistory/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)



func (r *PostgresRepository) MarkAsRead(ctx context.Context, userID uint, chapterID uint) (entity.ReadingHistory, error) {
	const op = richerror.Op("repository.postgres.readinghistory.MarkAsRead")

	query :=`
		INSERT INTO reading_history (user_id, chapter_id, read_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id,chapter_id)
		Do UPDATE SET read_at =NOW()
		RETURNING id, user_id, chapter_id, read_at
	`
	var history entity.ReadingHistory

	err := r.db.QueryRowContext(ctx, query,userID,chapterID).Scan(
		&history.ID,
		&history.UserID,
		&history.ChapterID,
		&history.ReadAt,
	)
	if err != nil {
		return entity.ReadingHistory{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to mark chapter as read")
	}

	return history, nil
}