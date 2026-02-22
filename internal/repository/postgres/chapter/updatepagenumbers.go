package chapter

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) UpdatePageNumbers(ctx context.Context, updates []entity.PageNumberUpdate) error {
	const op = richerror.Op("postgres.chapter.UpdatePageNumbers")

	if len(updates) == 0 {
		return nil
	}

	// Begin transaction — all updates succeed or all fail
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to begin transaction").
			WithKind(richerror.KindUnexpected)
	}
	defer tx.Rollback()

	for _, u := range updates {
		_, err := tx.ExecContext(ctx,
			`UPDATE chapter_pages SET page_number = $1 WHERE id = $2`,
			u.PageNumber, u.PageID,
		)
		if err != nil {
			logger.Error("failed to update page number",
				"page_id", u.PageID,
				"page_number", u.PageNumber,
				"error", err,
			)
			return richerror.New(op).
				WithErr(err).
				WithMessage("failed to update page number").
				WithKind(richerror.KindUnexpected)
		}
	}

	if err := tx.Commit(); err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to commit transaction").
			WithKind(richerror.KindUnexpected)
	}

	logger.Info("page numbers updated", "count", len(updates))

	return nil
}