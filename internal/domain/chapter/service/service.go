package service

import "github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/repository"


type Service struct {
	repo repository.Repository
}


func New(repo repository.Repository)Service{
 return Service{
	repo:repo,
 }
}