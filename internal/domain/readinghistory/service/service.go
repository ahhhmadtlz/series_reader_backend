package service

import (
	chapterRepository "github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/repository"
	readinghistoryRepository "github.com/ahhhmadtlz/series_reader_backend/internal/domain/readinghistory/repository"
	seriesRepository "github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/repository"
)

type Service struct {
	readinghistoryRepo readinghistoryRepository.Repository
	chapterRepo chapterRepository.Repository
	seriesRepo  seriesRepository.Repository
}

func New(readinghistoryRepo readinghistoryRepository.Repository, chapterRepo chapterRepository.Repository, seriesRepo seriesRepository.Repository) Service {
	return Service{
	  readinghistoryRepo:  readinghistoryRepo,
		chapterRepo: chapterRepo,
		seriesRepo:  seriesRepo,
	}
}