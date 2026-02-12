package chapter

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (r *PostgresRepository) GetByID(ctx context.Context,id uint) (*entity.Chapter,error){
	const op=richerror.Op("postgres.chapter.GetByID")

	logger.Debug("getting chapter by id ","chapter_id",id)

	query:=`
	 SELECT  id, series_id, chapter_number, title, created_at, updated_at
	 From chapters
	 WHERE id =$1
	`
	var  ch entity.Chapter

	err :=r.db.QueryRowContext(ctx,query,id).Scan(
		&ch.ID,
		&ch.SeriesID,
		&ch.ChapterNumber,
		&ch.Title,
		&ch.CreatedAt,
		&ch.UpdatedAt,
	)

	if err !=nil{
		if errors.Is(err,sql.ErrNoRows){
			return nil,richerror.New(op).WithErr(err).WithKind(richerror.KindNotFound)
		}
		logger.Error("failed to get chapter","error",err)
		return  nil,richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return &ch,nil
}


func (r *PostgresRepository) GetBySeriesID(ctx context.Context,seriesID uint)([]*entity.Chapter,error){
	const op = richerror.Op("postgres.chapter.GetBySeriesID")

	logger.Debug("getting chapters by series","series_id",seriesID)

	query := `
		SELECT id, series_id, chapter_number, title, created_at, updated_at
		FROM chapters
		WHERE series_id = $1
		ORDER BY chapter_number ASC
	`
	rows, err := r.db.QueryContext(ctx, query, seriesID)

	if err !=nil{
		logger.Error("failed to query chapters","error",err)
		return  nil,richerror.New(op).WithErr(err)
	}

	defer rows.Close()

	var chapters []*entity.Chapter
		for rows.Next() {
		var ch entity.Chapter
		if err := rows.Scan(
			&ch.ID,
			&ch.SeriesID,
			&ch.ChapterNumber,
			&ch.Title,
			&ch.CreatedAt,
			&ch.UpdatedAt,
		); err != nil {
			return nil, richerror.New(op).WithErr(err)
		}
		chapters = append(chapters, &ch)
	}

	return chapters, nil
}