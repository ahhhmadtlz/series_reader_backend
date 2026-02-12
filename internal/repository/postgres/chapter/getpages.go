package chapter

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) GetPagesByChapterID(ctx context.Context, chapterID uint) ([]entity.ChapterPage, error) {
	const op = richerror.Op("postgres.chapter.GetPagesByChapterID")

	query := `
		SELECT id, chapter_id, page_number, image_url, created_at
		FROM chapter_pages
		WHERE chapter_id = $1
		ORDER BY page_number ASC
	`

	rows, err := r.db.QueryContext(ctx, query, chapterID)
	if err != nil {
		return nil, richerror.New(op).WithErr(err)
	}
	defer rows.Close()

	var pages []entity.ChapterPage

	for rows.Next() {
		var p entity.ChapterPage
		if err := rows.Scan(
			&p.ID,
			&p.ChapterID,
			&p.PageNumber,
			&p.ImageURL,
			&p.CreatedAt,
		); err != nil {
			return nil, richerror.New(op).WithErr(err)
		}
		pages = append(pages, p)
	}

	return pages, nil
}



func (r *PostgresRepository) GetPageByNumber(ctx context.Context, chapterID uint, pageNumber int) (*entity.ChapterPage, error) {
	const op = richerror.Op("postgres.chapter.GetPageByNumber")

	query := `
		SELECT id, chapter_id, page_number, image_url, created_at
		FROM chapter_pages
		WHERE chapter_id = $1 AND page_number = $2
	`

	var p entity.ChapterPage

	err := r.db.QueryRowContext(ctx, query, chapterID, pageNumber).
		Scan(&p.ID, &p.ChapterID, &p.PageNumber, &p.ImageURL, &p.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, richerror.New(op).WithErr(err)
		}
		return nil, richerror.New(op).WithErr(err)
	}

	return &p, nil
}


func (r *PostgresRepository) GetChapterWithPages(ctx context.Context, chapterID uint) (*entity.Chapter, []entity.ChapterPage, error) {
	const op = richerror.Op("postgres.chapter.GetChapterWithPages")

	chapter, err := r.GetByID(ctx, chapterID)
	if err != nil {
		return nil, nil, richerror.New(op).WithErr(err)
	}

	pages, err := r.GetPagesByChapterID(ctx, chapterID)
	if err != nil {
		return nil, nil, richerror.New(op).WithErr(err)
	}

	return chapter, pages, nil
}
