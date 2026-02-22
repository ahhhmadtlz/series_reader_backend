package service

import (
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/repository"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/storage"
)


type Service struct {
	repo repository.Repository
	storage storage.Storage
}


func New(repo repository.Repository,storage storage.Storage)Service{
 return Service{
	repo:repo,
	storage: storage,
 }
}