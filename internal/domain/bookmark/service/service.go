package service

import (
	"context"

	BookmarkRepository "github.com/ahhhmadtlz/series_reader_backend/internal/domain/bookmark/repository"

	seriesEntity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/entity"
)



type SeriesRepository interface {
	GetByID(ctx context.Context, id uint) (seriesEntity.Series, error)
}

type Service struct {
	bookmarkRepo BookmarkRepository.Repository
	seriesRepo   SeriesRepository
}

func New(bookmarkRepo BookmarkRepository.Repository, seriesRepo SeriesRepository) Service {
	return Service{
		bookmarkRepo: bookmarkRepo,
		seriesRepo:   seriesRepo,
	}
}
