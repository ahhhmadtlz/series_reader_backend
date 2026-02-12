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

	logger.Debug("creating pages", "count", len(pages))

  query := `INSERT INTO chapter_pages (chapter_id, page_number, image_url) VALUES `
  values:=[]any{}
 
	for i , p :=range pages {
		if i >0 {
			query +=", "
		}
		query +=fmt.Sprintf("($%d, $%d, $%d)",i*3+1,i*3+2,i*3+3)
		values=append(values, p.ChapterID, p.PageNumber, p.ImageURL)
	}

	_,err :=r.db.ExecContext(ctx, query, values...)

	if err!=nil{
		logger.Error("failed to insert pages","error",err)
		return richerror.New(op).WithErr(err).WithMessage("failed to create pages").WithKind(richerror.KindUnexpected)
	}

	logger.Info("pages created","count",len(pages))

	return nil

}
