package chapter

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) CreatePages(ctx context.Context, pages []entity.ChapterPage) ([]entity.ChapterPage, error) {
	const op = richerror.Op("postgres.chapter.CreatePages")

	if len(pages) == 0 {
		return nil, nil
	}

	logger.Debug("creating pages", "count", len(pages))

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("failed to begin transaction").
			WithKind(richerror.KindUnexpected)
	}
	defer tx.Rollback()

	const query = `
		INSERT INTO chapter_pages (chapter_id, page_number, image_url, remote_path)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	created := make([]entity.ChapterPage, len(pages))
	for i, p := range pages {
		row := tx.QueryRowContext(ctx, query, p.ChapterID, p.PageNumber, p.ImageURL, p.RemotePath)
		if err := row.Scan(&p.ID); err != nil {
			logger.Error("failed to insert page", "error", err)
			return nil, richerror.New(op).
				WithErr(err).
				WithMessage("failed to create page").
				WithKind(richerror.KindUnexpected)
		}
		created[i] = p
	}

	if err := tx.Commit(); err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("failed to commit transaction").
			WithKind(richerror.KindUnexpected)
	}

	logger.Info("pages created", "count", len(created))

	return created, nil
}