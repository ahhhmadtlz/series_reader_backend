

package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/readinghistory/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) GetUserHistory(ctx context.Context, userID uint, limit int, offset int) ([]param.ReadingHistoryResponse, error) {
	const op = richerror.Op("service.readinghistory.GetUserHistory")

	
	histories, err := s.readinghistoryRepo.GetUserHistory(ctx, userID, limit, offset)
	if err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get user reading history")
	}

	return histories, nil
}