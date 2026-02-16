package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/bookmark/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/bookmark/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) CreateBookmark(ctx context.Context, userID uint, req param.CreateBookmarkRequest) (param.CreateBookmarkResponse, error) {
	const op = richerror.Op("bookmarkservice.CreateBookmark")

	logger.Info("Create bookmark request",
	"user_id",userID,
	"series_id",req.SeriesID,
	)

	series, err:=s.seriesRepo.GetByID(ctx,req.SeriesID)

	if err !=nil{
		logger.Error("Failed to get series",
			"series_id", req.SeriesID,
			"error", err.Error(),
		)
		return param.CreateBookmarkResponse{}, richerror.New(op).
			WithMessage("series not found").
			WithKind(richerror.KindNotFound).
			WithErr(err)
	}

	bookmark :=entity.Bookmark{
		UserID: userID,
		SeriesID: req.SeriesID,
	}

	createdBookmark, err := s.bookmarkRepo.Create(ctx, bookmark)
	if err != nil {
		logger.Error("Failed to create bookmark",
			"user_id", userID,
			"series_id", req.SeriesID,
			"error", err.Error(),
		)
		return param.CreateBookmarkResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to create bookmark").
			WithKind(richerror.KindUnexpected)
	}

	logger.Info("Bookmark created successfully",
		"bookmark_id", createdBookmark.ID,
		"user_id", userID,
		"series_id", req.SeriesID,
	)

	return param.CreateBookmarkResponse{
		Bookmark: param.BookmarkInfo{
			ID:        createdBookmark.ID,
			SeriesID:  createdBookmark.SeriesID,
			CreatedAt: createdBookmark.CreatedAt,
			Series: param.SeriesInfo{
				ID:            series.ID,
				Title:         series.Title,
				FullSlug:      series.FullSlug,
				CoverImageURL: series.CoverImageURL,
				Type:          series.Type,
				Status:        series.Status,
				Genres:        series.Genres,
			},
		},
	}, nil


}