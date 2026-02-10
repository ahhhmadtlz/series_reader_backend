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

func (s Service) GetByFullSlug(ctx context.Context, fullSlug string) (param.SeriesResponse, error) {
    const op = richerror.Op("service.series.GetByFullSlug")

    series, err := s.repo.GetByFullSlug(ctx, fullSlug)
    if err != nil {
        return param.SeriesResponse{}, richerror.New(op).WithErr(err).WithMessage("failed to get series by slug").WithKind(richerror.KindNotFound)
    }

    return toSeriesResponse(series), nil
}