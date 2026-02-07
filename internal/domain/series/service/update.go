package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/slugify"
)

func (s Service) Update(ctx context.Context, id uint, req param.UpdateSeriesRequest) (param.SeriesResponse, error) {
	const op = richerror.Op("service.series.Update")

	// Get existing series
	existingSeries, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return param.SeriesResponse{}, richerror.New(op).WithErr(err).WithMessage("failed to get existing series").WithKind(richerror.KindNotFound)
	}

	// Update only provided fields
	if req.Title != nil {
		existingSeries.Title = *req.Title
		// Regenerate slug if title changed
		newSlug := slugify.Make(*req.Title)
		if newSlug != existingSeries.Slug {
			exists, err := s.repo.IsSlugExists(ctx, newSlug)
			if err != nil {
				return param.SeriesResponse{}, richerror.New(op).WithErr(err).WithMessage("failed to check slug existence").WithKind(richerror.KindUnexpected)
			}
			if exists {
				newSlug = slugify.MakeUnique(*req.Title, func(testSlug string) bool {
					exists, _ := s.repo.IsSlugExists(ctx, testSlug)
					return exists
				})
			}
			existingSeries.Slug = newSlug
		}
	}

	if req.Description != nil {
		existingSeries.Description = *req.Description
	}
	if req.Author != nil {
		existingSeries.Author = *req.Author
	}
	if req.Artist != nil {
		existingSeries.Artist = *req.Artist
	}
	if req.Status != nil {
		existingSeries.Status = *req.Status
	}
	if req.Type != nil {
		existingSeries.Type = *req.Type
	}
	if req.Genres != nil {
		existingSeries.Genres = req.Genres
	}
	if req.AlternativeTitles != nil {
		existingSeries.AlternativeTitles = req.AlternativeTitles
	}
	if req.CoverImageURL != nil {
		existingSeries.CoverImageURL = *req.CoverImageURL
	}
	if req.PublicationYear != nil {
		existingSeries.PublicationYear = req.PublicationYear
	}
	if req.IsPublished != nil {
		existingSeries.IsPublished = *req.IsPublished
	}

	// Ensure slices are not nil
	if existingSeries.Genres == nil {
		existingSeries.Genres = []string{}
	}
	if existingSeries.AlternativeTitles == nil {
		existingSeries.AlternativeTitles = []string{}
	}

	updatedSeries, err := s.repo.Update(ctx, id, existingSeries)
	if err != nil {
		return param.SeriesResponse{}, richerror.New(op).WithErr(err).WithMessage("failed to update series").WithKind(richerror.KindUnexpected)
	}

	return toSeriesResponse(updatedSeries), nil
}
