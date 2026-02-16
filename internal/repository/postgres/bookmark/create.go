package bookmark

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/bookmark/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	"github.com/lib/pq"
)

func (r *PostgresRepository) Create(ctx context.Context, bookmark entity.Bookmark) (entity.Bookmark, error) {
	const op = richerror.Op("repository.postgres.bookmark.Create")

	query := `
		INSERT INTO bookmarks (user_id, series_id)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	var createdBookmark entity.Bookmark

	err := r.db.QueryRowContext(ctx, query, bookmark.UserID, bookmark.SeriesID).Scan(
		&createdBookmark.ID,
		&createdBookmark.CreatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// 23505 = unique_violation
			if pqErr.Code == "23505" {
				return entity.Bookmark{}, richerror.New(op).
					WithMessage("series already bookmarked").
					WithKind(richerror.KindInvalid)
			}
		}
		return entity.Bookmark{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to create bookmark")
	}

	createdBookmark.UserID = bookmark.UserID
	createdBookmark.SeriesID = bookmark.SeriesID

	return createdBookmark, nil
}