package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) DeleteImage(ctx context.Context, imageID uint) error {
	const op = richerror.Op("service.upload.DeleteImage")

	logger.Info("Delete image request", "image_id", imageID)

	// 1. Get image metadata (to know stored_path)
	img, err := s.uploadRepo.GetByID(ctx, imageID)
	if err != nil {
		logger.Error("Failed to get image",
			"image_id", imageID,
			"error", err.Error(),
		)
		return richerror.New(op).
			WithErr(err).
			WithMessage("image not found").
			WithKind(richerror.KindNotFound)
	}

	// 2. Delete DB record first (source of truth)
	err = s.uploadRepo.DeleteByID(ctx, imageID)
	if err != nil {
		logger.Error("Failed to delete image from database",
			"image_id", imageID,
			"error", err.Error(),
		)
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to delete image record").
			WithKind(richerror.KindUnexpected)
	}

	logger.Info("Image record deleted from database", "image_id", imageID)

	// 3. Delete physical file (best effort)
	err = s.storage.Delete(ctx, img.StoredPath)
	if err != nil {
		// Log but don't fail - file can be cleaned up async
		logger.Error("Failed to delete file after DB deletion",
			"image_id", imageID,
			"stored_path", img.StoredPath,
			"error", err.Error(),
		)
		// TODO: Add to cleanup queue for retry
	} else {
		logger.Info("File deleted successfully",
			"image_id", imageID,
			"stored_path", img.StoredPath,
		)
	}

	return nil
}