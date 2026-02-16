package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/bookmark/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) GetBookmarks(ctx context.Context, userID uint) (param.GetBookmarksResponse, error) {
	const op = richerror.Op("bookmarkservice.GetBookmarks")

	logger.Debug("Get bookmarks request", "user_id", userID)  // Changed to Debug

	bookmarks, err := s.bookmarkRepo.GetBookmarksWithSeriesByUserID(ctx, userID)
	if err != nil {
		logger.Error("Failed to get bookmarks",
			"user_id", userID,
			"error", err.Error(),
		)
		return param.GetBookmarksResponse{}, richerror.New(op).WithErr(err)
	}

	logger.Debug("Bookmarks retrieved",  // Changed to Debug
		"user_id", userID,
		"count", len(bookmarks),
	)

	return param.GetBookmarksResponse{
		Bookmarks: bookmarks,
		Total:     len(bookmarks),
	}, nil
}