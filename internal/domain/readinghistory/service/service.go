package service

import (
	"context"

	chapterEntity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/entity"
	readinghistoryRepository "github.com/ahhhmadtlz/series_reader_backend/internal/domain/readinghistory/repository"
	seriesEntity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/entity"
)

type ChapterRepository interface {
	GetByID(ctx context.Context, id uint) (*chapterEntity.Chapter, error)
	GetBySeriesID(ctx context.Context, seriesID uint) ([]*chapterEntity.Chapter, error)
}

type SeriesRepository interface {
	GetByID(ctx context.Context, id uint) (seriesEntity.Series, error)
	GetByFullSlug(ctx context.Context, slug string) (seriesEntity.Series, error)
}

type Service struct {
	readinghistoryRepo readinghistoryRepository.Repository
	chapterRepo        ChapterRepository
	seriesRepo         SeriesRepository
}

func New(
	readinghistoryRepo readinghistoryRepository.Repository,
	chapterRepo ChapterRepository,
	seriesRepo SeriesRepository,
) Service {
	return Service{
		readinghistoryRepo: readinghistoryRepo,
		chapterRepo:        chapterRepo,
		seriesRepo:         seriesRepo,
	}
}