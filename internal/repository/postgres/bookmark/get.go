package bookmark

import (
	"context"
	"encoding/json"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/bookmark/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) GetBookmarksWithSeriesByUserID(ctx context.Context, userID uint) ([]param.BookmarkInfo, error) {
	const op = richerror.Op("repository.postgres.bookmark.GetBookmarksWithSeriesByUserID")

	query := `
		SELECT  
		  b.id,
			b.series_id,
			b.created_at,
			s.id,
			s.title,
			s.full_slug,
			s.cover_image_url,
			s.type,
			s.status,
			s.genres
		FROM bookmarks b
		INNER JOIN series s ON s.id = b.series_id
		WHERE b.user_id = $1
		ORDER BY b.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get bookmarks with series")
	}
	defer rows.Close()

	bookmarks := []param.BookmarkInfo{}

	for rows.Next() {
		var bookmark param.BookmarkInfo
		var genresJSON []byte

		err := rows.Scan(
			&bookmark.ID,
			&bookmark.SeriesID,
			&bookmark.CreatedAt,
			&bookmark.Series.ID,
			&bookmark.Series.Title,
			&bookmark.Series.FullSlug,
			&bookmark.Series.CoverImageURL,
			&bookmark.Series.Type,
			&bookmark.Series.Status,
			&genresJSON,
		)
		if err != nil {
			return nil, richerror.New(op).
				WithErr(err).
				WithMessage("failed to scan bookmark")
		}

		// Unmarshal genres JSON
		if err := json.Unmarshal(genresJSON, &bookmark.Series.Genres); err != nil {
			return nil, richerror.New(op).
				WithErr(err).
				WithMessage("failed to unmarshal genres")
		}

		bookmarks = append(bookmarks, bookmark)
	}

	if err = rows.Err(); err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("error iterating bookmarks")
	}

	return bookmarks, nil
}