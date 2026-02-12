package chapter

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) Create(ctx context.Context,chapter *entity.Chapter)(*entity.Chapter,error){
 const op=richerror.Op("postgres.chapter.Create")

 logger.Debug("creating chapter","series_id",chapter.SeriesID,"chapter_number",chapter.ChapterNumber)

 query:=`
   INSERT INTO chapters (series_id, chapter_number, title)
	 VALUES($1, $2, $3)
	 RETURNING id, created_at, updated_at
 `

 err:=r.db.QueryRowContext(
	ctx,
	query,
	chapter.SeriesID,
	chapter.ChapterNumber,
	chapter.Title,
 ).Scan(
	&chapter.ID,
	&chapter.CreatedAt,
	&chapter.UpdatedAt,
 )

 if err !=nil {
	logger.Error("failed to create chapter","error",err,"series_id",chapter.SeriesID)
  return nil,richerror.New(op).WithErr(err)
 }

 logger.Info("chapter created","chapter_id",chapter.ID)
 return chapter,nil
}