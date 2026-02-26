package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) Delete(ctx context.Context, id uint) error {
	const op = richerror.Op("service.chapter.Delete")
	// 1. get all pages for this chapter
	pages, err := s.repo.GetPagesByChapterID(ctx, id)
	if err != nil {
		logger.Error("failed to get pages for chapter", "chapter_id", id, "error", err)
	}

	// 2. for each page, delete variant files, variant rows, source file
	for _, p := range pages {
		if err := s.ipSvc.DeletePageVariants(ctx, p.ID); err != nil {
			logger.Error("failed to delete variants for page", "page_id", p.ID, "error", err)
		}

		if p.RemotePath != "" {
			if err := s.storage.Delete(ctx, p.RemotePath); err != nil {
				logger.Error("failed to delete page file", "remote_path", p.RemotePath, "error", err)
			}
		}
	}


	// 3. delete the chapter (cascades page rows in DB)
	if err := s.repo.Delete(ctx, id); err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to delete chapter").
			WithKind(richerror.KindUnexpected)
	}

	logger.Info("chapter deleted", "chapter_id", id)

	return nil
}