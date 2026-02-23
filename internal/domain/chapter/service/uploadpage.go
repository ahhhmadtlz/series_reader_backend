package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/param"
	uploadentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/storage"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) UploadPage(ctx context.Context, req param.UploadPageParam) (param.ChapterPageResponse, error) {
	const op = richerror.Op("service.chapter.UploadPage")

	// 1. Verify chapter exists
	_, err := s.repo.GetByID(ctx, req.ChapterID)
	if err != nil {
		return param.ChapterPageResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("chapter not found").
			WithKind(richerror.KindNotFound)
	}

	// 2. Validate page number is non-negative
	if req.PageNumber < 0 {
		return param.ChapterPageResponse{}, richerror.New(op).
			WithMessage("page number must be non-negative").
			WithKind(richerror.KindInvalid)
	}

	// 3. Save file to storage
	result, err := s.storage.Save(ctx, storage.SaveRequest{
		File:     req.File,
		Filename: req.Header.Filename,
		OwnerID:  req.ChapterID,
		Kind:     uploadentity.ImageKindChapterPage,
		MimeType: req.Header.Header.Get("Content-Type"),
		Size:     req.Header.Size,
	})
	if err != nil {
		return param.ChapterPageResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to save page file").
			WithKind(richerror.KindUnexpected)
	}

	// 4. Save record to DB
	pages := []entity.ChapterPage{
		{
			ChapterID:  req.ChapterID,
			PageNumber: req.PageNumber,
			ImageURL:   result.URL,
			RemotePath: result.StoredPath,
		},
	}

	if _,err := s.repo.CreatePages(ctx, pages); err != nil {
		// Rollback: delete file from storage
		_ = s.storage.Delete(ctx, result.StoredPath)
		return param.ChapterPageResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to save page record").
			WithKind(richerror.KindUnexpected)
	}

	logger.Info("page uploaded", "chapter_id", req.ChapterID, "page_number", req.PageNumber)

	return param.ChapterPageResponse{
		PageNumber: req.PageNumber,
		ImageURL:   result.URL,
	}, nil
}