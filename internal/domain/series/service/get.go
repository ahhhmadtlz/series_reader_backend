package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) GetByID(ctx context.Context, id uint) (param.SeriesResponse, error) {
	const op = richerror.Op("service.series.GetByID")

	series, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return param.SeriesResponse{}, richerror.New(op).WithErr(err).WithMessage("failed to get series").WithKind(richerror.KindNotFound)
	}

	return toSeriesResponse(series), nil
}

func (s Service) GetBySlug(ctx context.Context, slug string) (param.SeriesResponse, error) {
	const op = richerror.Op("service.series.GetBySlug")

	series, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return param.SeriesResponse{}, richerror.New(op).WithErr(err).WithMessage("failed to get series").WithKind(richerror.KindNotFound)
	}

	return toSeriesResponse(series), nil
}