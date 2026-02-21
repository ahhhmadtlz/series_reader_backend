package service

import (
	seriesrepo "github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/repository"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/repository"
	userrepo "github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/repository"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/storage"
)


type Service struct {
	uploadRepo repository.Repository
	userRepo    userrepo.Repository
	seriesRepo  seriesrepo.Repository
	storage   storage.Storage
}

func New(
	uploadRepo repository.Repository,
	userRepo userrepo.Repository,
	seriesRepo seriesrepo.Repository,
	storage storage.Storage,
) Service {
	return Service{
		uploadRepo: uploadRepo,
		userRepo: userRepo,
		seriesRepo: seriesRepo,
		storage: storage,
	}
}