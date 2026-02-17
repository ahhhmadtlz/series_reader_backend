package repository

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/readinghistory/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/readinghistory/param"
)

type Repository interface {
	MarkAsRead(ctx context.Context, userID uint, chapterID uint) (entity.ReadingHistory, error)
	IsChapterRead(ctx context.Context, userID uint, chapterID uint) (bool, error)
  GetUserHistory(ctx context.Context, userID uint, limit int, offset int) ([]param.ReadingHistoryResponse, error)
	GetSeriesProgress(ctx context.Context, userID uint, seriesID uint) ([]entity.ReadingHistory, error)
	UnmarkAsRead(ctx context.Context, userID uint, chapterID uint) error
	GetTotalReadCount(ctx context.Context, userID uint) (int, error)

}