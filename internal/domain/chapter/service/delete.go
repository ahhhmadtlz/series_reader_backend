package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) Delete(ctx context.Context, id uint) error {
	const op = richerror.Op("service.chapter.Delete")

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to delete chapter").
			WithKind(richerror.KindUnexpected)
	}

	return nil
}
