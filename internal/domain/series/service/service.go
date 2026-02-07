package service

import (
	"context"
	"time"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/param"
)

type Repository interface {
	Create(ctx context.Context, series entity.Series) (entity.Series, error)
	GetByID(ctx context.Context, id uint) (entity.Series, error)
	GetBySlug(ctx context.Context, slug string) (entity.Series, error)
	GetList(ctx context.Context, req param.GetListRequest) ([]entity.Series, int, error)
	Update(ctx context.Context, id uint, series entity.Series) (entity.Series, error)
	Delete(ctx context.Context, id uint) error
	IncrementViewCount(ctx context.Context, id uint) error
	IsSlugExists(ctx context.Context, slug string) (bool, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) Service {
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