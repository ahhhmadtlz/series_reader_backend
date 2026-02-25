package service

import (
	"context"

	ipparam "github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/storage"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) UploadChapterThumbnail(ctx context.Context, req param.UploadChapterThumbnailRequest) (param.UploadChapterThumbnailResponse, error) {
	const op = richerror.Op("service.upload.UploadChapterThumbnail")

	logger.Info("Upload chapter thumbnail request",
		"chapter_id", req.ChapterID,
		"user_id", req.UserID,
		"filename", req.Header.Filename,
		"size", req.Header.Size,
	)

	// 1. Check chapter exists
	_, err := s.chapterRepo.GetByID(ctx, req.ChapterID)
	if err != nil {
		return param.UploadChapterThumbnailResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("chapter not found").
			WithKind(richerror.KindNotFound)
	}

	// 2. Delete old thumbnail if exists (best effort)
	oldThumb, err := s.uploadRepo.GetLatestByOwner(ctx, req.ChapterID, entity.ImageKindChapterThumbnail)
	if err == nil {
		_ = s.uploadRepo.DeleteByID(ctx, oldThumb.ID)
		if err := s.storage.Delete(ctx, oldThumb.StoredPath); err != nil {
			logger.Error("failed to delete old thumbnail file",
				"chapter_id", req.ChapterID,
				"stored_path", oldThumb.StoredPath,
				"error", err,
			)
		}
	}

	// 3. Save to storage
	result, err := s.storage.Save(ctx, storage.SaveRequest{
		File:     req.File,
		Filename: req.Header.Filename,
		OwnerID:  req.ChapterID,
		Kind:     entity.ImageKindChapterThumbnail,
		MimeType: req.Header.Header.Get("Content-Type"),
		Size:     req.Header.Size,
	})
	if err != nil {
		return param.UploadChapterThumbnailResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to save thumbnail file").
			WithKind(richerror.KindUnexpected)
	}

	// 4. Save record to DB
	thumbImg := entity.UploadedImage{
		OwnerID:    req.ChapterID,
		Kind:       entity.ImageKindChapterThumbnail,
		Filename:   req.Header.Filename,
		StoredPath: result.StoredPath,
		URL:        result.URL,
		MimeType:   req.Header.Header.Get("Content-Type"),
		SizeBytes:  req.Header.Size,
	}

	savedImg, err := s.uploadRepo.Save(ctx, thumbImg)
	if err != nil {
		_ = s.storage.Delete(ctx, result.StoredPath)
		return param.UploadChapterThumbnailResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to save thumbnail record").
			WithKind(richerror.KindUnexpected)
	}

	logger.Info("Chapter thumbnail uploaded successfully",
		"chapter_id", req.ChapterID,
		"thumbnail_url", savedImg.URL,
	)

	// 5. Enqueue image processing job (fire and forget)
	if err := s.jobQueue.Enqueue(ctx, ipparam.ProcessImageArgs{
		OwnerID:    req.ChapterID,
		RemotePath: savedImg.StoredPath,
		ImageKind:  entity.ImageKindChapterThumbnail,
	}); err != nil {
		logger.Error("failed to enqueue thumbnail processing job",
			"chapter_id", req.ChapterID,
			"error", err,
		)
	}

	return param.UploadChapterThumbnailResponse{
		ThumbnailURL: savedImg.URL,
	}, nil
}