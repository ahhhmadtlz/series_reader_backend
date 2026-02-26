package service

import (
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/repository"
	iprepo "github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/repository"
	ipservice "github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/service"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/storage"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/worker"
)

type Service struct {
	repo        repository.Repository
	storage     storage.Storage
	jobQueue    worker.JobQueue
	variantRepo iprepo.Repository
	ipSvc       ipservice.Service
}

func New(
	repo repository.Repository,
	store storage.Storage,
	jobQueue worker.JobQueue,
	variantRepo iprepo.Repository,
	ipSvc ipservice.Service,
) Service {
	return Service{
		repo:        repo,
		storage:     store,
		jobQueue:    jobQueue,
		variantRepo: variantRepo,
		ipSvc:       ipSvc,
	}
}