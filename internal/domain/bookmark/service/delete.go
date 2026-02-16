package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/bookmark/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) DeleteBookmark(ctx context.Context, userID uint, seriesID uint) (param.DeleteBookmarkResponse, error) {
	const op = richerror.Op("bookmarkservice.DeleteBookmark")

	logger.Debug("Delete bookmark", "user_id", userID, "series_id", seriesID)

	err := s.bookmarkRepo.Delete(ctx, userID, seriesID)
	if err != nil {
		logger.Error("Failed to delete bookmark",
			"user_id", userID,
			"series_id", seriesID,
			"error", err.Error(),
		)
		return param.DeleteBookmarkResponse{}, richerror.New(op).WithErr(err)
	}

	logger.Debug("Bookmark deleted", "user_id", userID, "series_id", seriesID)

	return param.DeleteBookmarkResponse{
		Message: "bookmark deleted successfully",
	}, nil
}