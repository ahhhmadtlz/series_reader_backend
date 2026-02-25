package service

import (
	seriesrepo "github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/repository"
	chapterrepo "github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/repository"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/repository"
	userrepo "github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/repository"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/storage"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/worker"
)

type Service struct {
	uploadRepo  repository.Repository
	userRepo    userrepo.Repository
	seriesRepo  seriesrepo.Repository
	chapterRepo chapterrepo.Repository
	storage     storage.Storage
	jobQueue    worker.JobQueue
}

func New(
	uploadRepo repository.Repository,
	userRepo userrepo.Repository,
	seriesRepo seriesrepo.Repository,
	chapterRepo chapterrepo.Repository,
	storage storage.Storage,
	jobQueue worker.JobQueue,
) Service {
	return Service{
		uploadRepo:  uploadRepo,
		userRepo:    userRepo,
		seriesRepo:  seriesRepo,
		chapterRepo: chapterRepo,
		storage:     storage,
		jobQueue:    jobQueue,
	}
}