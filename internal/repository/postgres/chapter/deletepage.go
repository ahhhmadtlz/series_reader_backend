package chapter

import (
	"context"
	"database/sql"

	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

// DeletePage deletes a specific page by ID
func (r *PostgresRepository) DeletePage(ctx context.Context, pageID uint) error {
	const op = richerror.Op("postgres.chapter.DeletePage")

	logger.Info("Deleting page", "page_id", pageID)

	query := `DELETE FROM chapter_pages WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, pageID)
	if err != nil {
		logger.Error("Failed to delete page", "error", err, "page_id", pageID)
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to delete page").
			WithKind(richerror.KindUnexpected)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to get rows affected").
			WithKind(richerror.KindUnexpected)
	}

	if rowsAffected == 0 {
		return richerror.New(op).
			WithErr(sql.ErrNoRows).
			WithMessage("page not found").
			WithKind(richerror.KindNotFound)
	}

	logger.Info("Page deleted successfully", "page_id", pageID)
	return nil
}