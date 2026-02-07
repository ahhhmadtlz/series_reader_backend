package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) IncrementViewCount(ctx context.Context, id uint) error {
	const op = richerror.Op("service.series.IncrementViewCount")

	err := s.repo.IncrementViewCount(ctx, id)
	if err != nil {
		return richerror.New(op).WithErr(err).WithMessage("failed to increment view count").WithKind(richerror.KindUnexpected)
	}

	return nil
}
