package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)


func (s Service) UnmarkAsRead(ctx context.Context,chapterID uint ,userID uint)error {
	const op= richerror.Op("service.readinghistory.UnmarkAsRead")

	_,err:=s.chapterRepo.GetByID(ctx,chapterID)
	if err !=nil{
		return richerror.New(op).
		WithErr(err).
		WithMessage("chapter not found").
		WithKind(richerror.KindNotFound)
	}
	err = s.readinghistoryRepo.UnmarkAsRead(ctx , userID,chapterID)

	if err !=nil{
		return richerror.New(op).
		 WithErr(err).
		 WithMessage("failed to unmaark chapter as read")
	}
	
	return nil
}