package service

import (
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/repository"
	iprepo "github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/repository"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/storage"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/worker"
)

type Service struct {
	repo        repository.Repository
	storage     storage.Storage
	jobQueue    worker.JobQueue
	variantRepo iprepo.Repository
}

func New(
	repo repository.Repository,
	storage storage.Storage,
	jobQueue worker.JobQueue,
	variantRepo iprepo.Repository,
) Service {
	return Service{
		repo:        repo,
		storage:     storage,
		jobQueue:    jobQueue,
		variantRepo: variantRepo,
	}
}