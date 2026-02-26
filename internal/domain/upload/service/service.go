package service

import (
	"context"

	chapterEntity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/entity"
	seriesEntity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/repository"
	userEntity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/storage"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/worker"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, userID uint) (userEntity.User, error)
	UpdateUser(ctx context.Context, user userEntity.User) (userEntity.User, error)
}

type SeriesRepository interface {
	GetByID(ctx context.Context, id uint) (seriesEntity.Series, error)
	Update(ctx context.Context, id uint, series seriesEntity.Series) (seriesEntity.Series, error)
}

type ChapterRepository interface {
	GetByID(ctx context.Context, id uint) (*chapterEntity.Chapter, error)
}

type Service struct {
	uploadRepo  repository.Repository
	userRepo    UserRepository
	seriesRepo  SeriesRepository
	chapterRepo ChapterRepository
	storage     storage.Storage
	jobQueue    worker.JobQueue
}

func New(
	uploadRepo repository.Repository,
	userRepo UserRepository,
	seriesRepo SeriesRepository,
	chapterRepo ChapterRepository,
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