package service

import (
	"context"
	"fmt"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/param"
	ipparam "github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/param"
	sharedentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/shared/entity"
	uploadentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/storage"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) BulkUploadPages(ctx context.Context, req param.BulkUploadParam) ([]param.ChapterPageResponse, error) {
	const op = richerror.Op("service.chapter.BulkUploadPages")

	// 1. verify chapter exists
	_, err := s.repo.GetByID(ctx, req.ChapterID)
	if err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("chapter not found").
			WithKind(richerror.KindNotFound)
	}

	if len(req.Files) == 0 {
		return nil, richerror.New(op).
			WithMessage("no files provided").
			WithKind(richerror.KindInvalid)
	}

	// 2. check existing pages
	existing, err := s.repo.GetPagesByChapterID(ctx, req.ChapterID)
	if err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get existing pages").
			WithKind(richerror.KindUnexpected)
	}

	// 3. if pages exist and force is false, block with 409
	if len(existing) > 0 && !req.Force {
		return nil, richerror.New(op).
			WithMessage("chapter already has pages, send force=true to replace them").
			WithKind(richerror.KindConflict)
	}

	// 4. force=true is restricted to admin only
	// TODO: when series_admin collaboration role is added, also allow CallerCollabRole == SeriesAdminCollabRole here
	if req.Force && req.CallerRole != sharedentity.AdminRole {
		return nil, richerror.New(op).
			WithMessage("only admins can replace existing chapter pages").
			WithKind(richerror.KindForbidden)
	}

	// 5. delete existing pages, their variant files, and source files
	if len(existing) > 0 {
		for _, p := range existing {
			if err := s.ipSvc.DeletePageVariants(ctx, p.ID); err != nil {
				logger.Error("failed to delete variants for page", "page_id", p.ID, "error", err)
			}

			if err := s.repo.DeletePage(ctx, p.ID); err != nil {
				logger.Error("failed to delete existing page", "page_id", p.ID, "error", err)
			}

			if p.RemotePath != "" {
				if err := s.storage.Delete(ctx, p.RemotePath); err != nil {
					logger.Error("failed to delete existing page file", "remote_path", p.RemotePath, "error", err)
				}
			}
		}
	}

	// 6. upload each new file
	var pages []entity.ChapterPage
	var storedPaths []string // for rollback

	for i, fileHeader := range req.Files {
		file, err := fileHeader.Open()
		if err != nil {
			for _, path := range storedPaths {
				_ = s.storage.Delete(ctx, path)
			}
			return nil, richerror.New(op).
				WithErr(err).
				WithMessage(fmt.Sprintf("failed to open file %d", i+1)).
				WithKind(richerror.KindUnexpected)
		}

		result, err := s.storage.Save(ctx, storage.SaveRequest{
			File:     file,
			Filename: fileHeader.Filename,
			OwnerID:  req.ChapterID,
			Kind:     uploadentity.ImageKindChapterPage,
			MimeType: fileHeader.Header.Get("Content-Type"),
			Size:     fileHeader.Size,
		})
		file.Close()
		if err != nil {
			for _, path := range storedPaths {
				_ = s.storage.Delete(ctx, path)
			}
			return nil, richerror.New(op).
				WithErr(err).
				WithMessage(fmt.Sprintf("failed to save file %d", i+1)).
				WithKind(richerror.KindUnexpected)
		}

		storedPaths = append(storedPaths, result.StoredPath)
		pages = append(pages, entity.ChapterPage{
			ChapterID:  req.ChapterID,
			PageNumber: i + 1,
			ImageURL:   result.URL,
			RemotePath: result.StoredPath,
		})
	}

	// 7. save all records in one DB call
	created, err := s.repo.CreatePages(ctx, pages)
	if err != nil {
		for _, path := range storedPaths {
			_ = s.storage.Delete(ctx, path)
		}
		return nil, richerror.New(op).
			WithMessage("failed to save pages to database").
			WithKind(richerror.KindUnexpected)
	}

	logger.Info("bulk pages uploaded", "chapter_id", req.ChapterID, "count", len(pages))

	// 8. enqueue image processing job per page (fire and forget)
	for _, p := range created {
		if err := s.jobQueue.Enqueue(ctx, ipparam.ProcessImageArgs{
			PageID:     p.ID,
			OwnerID:    p.ChapterID,
			RemotePath: p.RemotePath,
			ChapterID:  p.ChapterID,
			ImageKind:  uploadentity.ImageKindChapterPage,
		}); err != nil {
			logger.Error("failed to enqueue image processing job", "page_id", p.ID, "error", err)
		}
	}

	responses := make([]param.ChapterPageResponse, len(created))
	for i, p := range created {
		responses[i] = param.ChapterPageResponse{
			ID:         p.ID,
			PageNumber: p.PageNumber,
			ImageURL:   p.ImageURL,
		}
	}

	return responses, nil
}