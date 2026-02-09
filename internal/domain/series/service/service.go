package service

import (
	"time"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/repository"
)


type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) Service {
	return Service{
		repo: repo,
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