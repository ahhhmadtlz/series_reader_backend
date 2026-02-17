package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/readinghistory/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) MarkAsRead(ctx context.Context, req param.MarkAsReadRequest, userID uint) (param.ReadingHistoryResponse, error) {
	const op = richerror.Op("service.readinghistory.MarkAsRead")
	chapter, err :=s.chapterRepo.GetByID(ctx, req.ChapterID)

	if err !=nil{
		return param.ReadingHistoryResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("chapter not found").
			WithKind(richerror.KindNotFound)
	}

	history,err :=s.readinghistoryRepo.MarkAsRead(ctx,userID,req.ChapterID)
	if err !=nil{
		return  param.ReadingHistoryResponse{},richerror.New(op).
		 WithErr(err).
		 WithMessage("failed to mark chapter as read")
	}

	series, err :=s.seriesRepo.GetByID(ctx,chapter.SeriesID)
	if err != nil {
		return param.ReadingHistoryResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get series details")
	}

	response := param.ReadingHistoryResponse{
		ID:            history.ID,
		UserID:        history.UserID,
		ChapterID:     history.ChapterID,
		SeriesID:      chapter.SeriesID,
		ChapterNumber: chapter.ChapterNumber,
		SeriesTitle:   series.Title,
		ReadAt:        history.ReadAt,
	}

	return response, nil
}