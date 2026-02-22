package chapter

import (
	"context"
	"fmt"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) CreatePages(ctx context.Context, pages []entity.ChapterPage) error {
	const op = richerror.Op("postgres.chapter.CreatePages")

	if len(pages) == 0 {
		return nil
	}

	logger.Debug("creating pages", "count", len(pages))

	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to begin transaction").
			WithKind(richerror.KindUnexpected)
	}
	defer tx.Rollback()

	query := `INSERT INTO chapter_pages (chapter_id, page_number, image_url, remote_path) VALUES `
	values := []any{}

	for i, p := range pages {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
		values = append(values, p.ChapterID, p.PageNumber, p.ImageURL, p.RemotePath)
	}

	_, err = tx.ExecContext(ctx, query, values...)
	if err != nil {
		logger.Error("failed to insert pages", "error", err)
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to create pages").
			WithKind(richerror.KindUnexpected)
	}

	if err := tx.Commit(); err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to commit transaction").
			WithKind(richerror.KindUnexpected)
	}

	logger.Info("pages created", "count", len(pages))

	return nil
}