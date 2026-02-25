package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/param"
	ipparam "github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/param"
	uploadentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/storage"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) UploadPage(ctx context.Context, req param.UploadPageParam) (param.ChapterPageResponse, error) {
	const op = richerror.Op("service.chapter.UploadPage")

	// 1. verify chapter exists
	_, err := s.repo.GetByID(ctx, req.ChapterID)
	if err != nil {
		return param.ChapterPageResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("chapter not found").
			WithKind(richerror.KindNotFound)
	}

	if req.PageNumber < 0 {
		return param.ChapterPageResponse{}, richerror.New(op).
			WithMessage("page number must be non-negative").
			WithKind(richerror.KindInvalid)
	}

	// 2. if a page already exists at this page number, clean it up first
	existing, err := s.repo.GetPageByNumber(ctx, req.ChapterID, req.PageNumber)
	if err == nil && existing != nil {
		// delete variant files and rows
		variants, err := s.variantRepo.GetVariantsByPageID(ctx, existing.ID)
		if err != nil {
			logger.Error("failed to get variants for existing page", "page_id", existing.ID, "error", err)
		}
		for _, v := range variants {
			if v.RemotePath != "" {
				if err := s.storage.Delete(ctx, v.RemotePath); err != nil {
					logger.Error("failed to delete variant file", "remote_path", v.RemotePath, "error", err)
				}
			}
		}
		if err := s.variantRepo.DeleteVariantsByPageID(ctx, existing.ID); err != nil {
			logger.Error("failed to delete variant rows", "page_id", existing.ID, "error", err)
		}
		// delete page row
		if err := s.repo.DeletePage(ctx, existing.ID); err != nil {
			logger.Error("failed to delete existing page row", "page_id", existing.ID, "error", err)
		}
		// delete source file
		if existing.RemotePath != "" {
			if err := s.storage.Delete(ctx, existing.RemotePath); err != nil {
				logger.Error("failed to delete existing page file", "remote_path", existing.RemotePath, "error", err)
			}
		}
	}

	// 3. save new file to storage
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

	// 4. save record to DB
	created, err := s.repo.CreatePages(ctx, []entity.ChapterPage{
		{
			ChapterID:  req.ChapterID,
			PageNumber: req.PageNumber,
			ImageURL:   result.URL,
			RemotePath: result.StoredPath,
		},
	})
	if err != nil {
		_ = s.storage.Delete(ctx, result.StoredPath)
		return param.ChapterPageResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to save page record").
			WithKind(richerror.KindUnexpected)
	}

	// 5. enqueue image processing job (fire and forget)
	if err := s.jobQueue.Enqueue(ctx, ipparam.ProcessImageArgs{
		PageID:     created[0].ID,
		OwnerID:    created[0].ChapterID,
		RemotePath: created[0].RemotePath,
		ChapterID:  created[0].ChapterID,
		ImageKind:  uploadentity.ImageKindChapterPage,
	}); err != nil {
		logger.Error("failed to enqueue image processing job", "page_id", created[0].ID, "error", err)
	}

	logger.Info("page uploaded", "chapter_id", req.ChapterID, "page_number", req.PageNumber)

	return param.ChapterPageResponse{
		ID:         created[0].ID,
		PageNumber: req.PageNumber,
		ImageURL:   result.URL,
	}, nil
}