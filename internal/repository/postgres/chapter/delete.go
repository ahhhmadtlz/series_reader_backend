package chapter

import (
	"context"
	"database/sql"

	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) Delete(ctx context.Context, id uint) error {
	const op = richerror.Op("postgres.chapter.Delete")

	logger.Info("deleting chapter", "chapter_id", id)

	query := `DELETE FROM chapters WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		logger.Error("failed to delete chapter", "error", err)
		return richerror.New(op).WithErr(err)
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return richerror.New(op).WithErr(sql.ErrNoRows).WithKind(richerror.KindNotFound)
	}

	return nil
}
