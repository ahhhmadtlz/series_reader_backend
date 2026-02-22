package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) DeletePage(ctx context.Context, chapterID uint, pageNumber int) error {
	const op = richerror.Op("service.chapter.DeletePage")

	// 1. Get the page to retrieve its remote path
	page, err := s.repo.GetPageByNumber(ctx, chapterID, pageNumber)
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("page not found").
			WithKind(richerror.KindNotFound)
	}

	// 2. Delete from DB first
	if err := s.repo.DeletePage(ctx, page.ID); err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to delete page").
			WithKind(richerror.KindUnexpected)
	}

	// 3. Delete file from storage (best effort — log but don't fail)
	if page.RemotePath != "" {
		if err := s.storage.Delete(ctx, page.RemotePath); err != nil {
			logger.Error("failed to delete page file from storage",
				"chapter_id", chapterID,
				"page_number", pageNumber,
				"remote_path", page.RemotePath,
				"error", err,
			)
		}
	}

	logger.Info("page deleted", "chapter_id", chapterID, "page_number", pageNumber)

	return nil
}