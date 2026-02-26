package service

import (
	"time"

	chapterrepo "github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/repository"
	ipservice "github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/service"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/repository"
	uploadrepository "github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/repository"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/storage"
)

type Service struct {
	repo        repository.Repository
	storage     storage.Storage
	uploadRepo  uploadrepository.Repository
	chapterRepo chapterrepo.Repository
	ipSvc       ipservice.Service
}

func New(
	repo repository.Repository,
	store storage.Storage,
	uploadRepo uploadrepository.Repository,
	chapterRepo chapterrepo.Repository,
	ipSvc ipservice.Service,
) Service {
	return Service{
		repo:        repo,
		storage:     store,
		uploadRepo:  uploadRepo,
		chapterRepo: chapterRepo,
		ipSvc:       ipSvc,
	}
}

func toSeriesResponse(series entity.Series) param.SeriesResponse {
	var createdAt, updatedAt string

	if !series.CreatedAt.IsZero() {
		createdAt = series.CreatedAt.Format(time.RFC3339)
	}

	if !series.UpdatedAt.IsZero() {
		updatedAt = series.UpdatedAt.Format(time.RFC3339)
	}

	return param.SeriesResponse{
		ID:                series.ID,
		Title:             series.Title,
		Slug:              series.Slug,
		SlugID:            series.SlugID,
		FullSlug:          series.FullSlug,
		Description:       series.Description,
		Author:            series.Author,
		Artist:            series.Artist,
		Status:            series.Status,
		Type:              series.Type,
		Genres:            series.Genres,
		AlternativeTitles: series.AlternativeTitles,
		CoverImageURL:     series.CoverImageURL,
		PublicationYear:   series.PublicationYear,
		ViewCount:         series.ViewCount,
		Rating:            series.Rating,
		IsPublished:       series.IsPublished,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}