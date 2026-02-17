package readinghistory

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/readinghistory/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/readinghistory/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

// GetUserHistory - Returns fully hydrated reading history with JOINs (no N+1)
func (r *PostgresRepository) GetUserHistory(ctx context.Context, userID uint, limit int, offset int) ([]param.ReadingHistoryResponse, error) {
	const op = richerror.Op("repository.postgres.readinghistory.GetUserHistory")

	query := `
		SELECT 
			rh.id,
			rh.user_id,
			rh.chapter_id,
			rh.read_at,
			c.series_id,
			c.chapter_number,
			s.title
		FROM reading_history rh
		JOIN chapters c ON rh.chapter_id = c.id
		JOIN series s ON c.series_id = s.id
		WHERE rh.user_id = $1
		ORDER BY rh.read_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get user reading history")
	}
	defer rows.Close()

	var histories []param.ReadingHistoryResponse
	for rows.Next() {
		var h param.ReadingHistoryResponse
		err := rows.Scan(
			&h.ID,
			&h.UserID,
			&h.ChapterID,
			&h.ReadAt,
			&h.SeriesID,
			&h.ChapterNumber,
			&h.SeriesTitle,
		)
		if err != nil {
			return nil, richerror.New(op).
				WithErr(err).
				WithMessage("failed to scan reading history")
		}
		histories = append(histories, h)
	}

	if err = rows.Err(); err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("error iterating reading history rows")
	}

	return histories, nil
}

func (r *PostgresRepository) GetSeriesProgress(ctx context.Context, userID uint, seriesID uint) ([]entity.ReadingHistory, error) {

	const op = richerror.Op("repository.postgres.readinghistory.GetSeriesProgress")

		query := `
		SELECT rh.id, rh.user_id, rh.chapter_id, rh.read_at
		FROM reading_history rh
		INNER JOIN chapters c ON rh.chapter_id = c.id
		WHERE rh.user_id = $1 AND c.series_id = $2
		ORDER BY c.chapter_number ASC
	 `
	  rows, err := r.db.QueryContext(ctx, query, userID, seriesID)
		if err != nil {
			return nil, richerror.New(op).
				WithErr(err).
				WithMessage("failed to get series progress")
		}
		defer rows.Close()

		var histories []entity.ReadingHistory
		for rows.Next() {
			var h entity.ReadingHistory
			err := rows.Scan(&h.ID, &h.UserID, &h.ChapterID, &h.ReadAt)
			if err != nil {
				return nil, richerror.New(op).
					WithErr(err).
					WithMessage("failed to scan series progress")
			}
			histories = append(histories, h)
		}

		if err = rows.Err(); err != nil {
			return nil, richerror.New(op).
				WithErr(err).
				WithMessage("error iterating series progress rows")
		}

		return histories, nil

}
func (r *PostgresRepository) GetTotalReadCount(ctx context.Context, userID uint) (int, error) {
	const op = richerror.Op("repository.postgres.readinghistory.GetTotalReadCount")

	query := `
		SELECT COUNT(*) 
		FROM reading_history
		WHERE user_id = $1
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)

	if err != nil {
		return 0, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get total read count")
	}

	return count, nil
}
