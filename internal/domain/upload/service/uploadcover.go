package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/storage"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) UploadCover(ctx context.Context, req param.UploadCoverRequest) (param.UploadCoverResponse, error) {
	const op = richerror.Op("service.upload.UploadCover")

	logger.Info("Upload cover request",
		"series_id", req.SeriesID,
		"user_id", req.UserID,
		"filename", req.Header.Filename,
		"size", req.Header.Size,
	)

	// 1. Check if series exists
	series, err := s.seriesRepo.GetByID(ctx, req.SeriesID)
	if err != nil {
		logger.Error("Failed to get series",
			"series_id", req.SeriesID,
			"error", err.Error(),
		)
		return param.UploadCoverResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("series not found").
			WithKind(richerror.KindNotFound)
	}

	logger.Info("Series found",
		"series_id", req.SeriesID,
		"title", series.Title,
	)

	// 2. Delete old cover if exists (best effort)
	oldCover, err := s.uploadRepo.GetLatestByOwner(ctx, req.SeriesID, entity.ImageKindCover)
	if err == nil {
		// Old cover exists, delete it
		logger.Info("Deleting old cover",
			"series_id", req.SeriesID,
			"old_cover_id", oldCover.ID,
			"old_path", oldCover.StoredPath,
		)

		// Delete from DB first
		_ = s.uploadRepo.DeleteByID(ctx, oldCover.ID)

		// Delete from storage (best effort, log if fails)
		if err := s.storage.Delete(ctx, oldCover.StoredPath); err != nil {
			logger.Error("Failed to delete old cover file",
				"series_id", req.SeriesID,
				"stored_path", oldCover.StoredPath,
				"error", err.Error(),
			)
			// Don't fail the upload if old file deletion fails
		}
	}

	// 3. Save new cover to storage
	saveReq := storage.SaveRequest{
		File:     req.File,
		Filename: req.Header.Filename,
		OwnerID:  req.SeriesID,
		Kind:     entity.ImageKindCover,
		MimeType: req.Header.Header.Get("Content-Type"),
		Size:     req.Header.Size,
	}

	result, err := s.storage.Save(ctx, saveReq)
	if err != nil {
		logger.Error("Failed to save cover to storage",
			"series_id", req.SeriesID,
			"error", err.Error(),
		)
		return param.UploadCoverResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to save cover file").
			WithKind(richerror.KindUnexpected)
	}

	logger.Info("Cover saved to storage",
		"series_id", req.SeriesID,
		"stored_path", result.StoredPath,
		"url", result.URL,
	)

	// 4. Save cover record to DB
	coverImg := entity.UploadedImage{
		OwnerID:    req.SeriesID,
		Kind:       entity.ImageKindCover,
		Filename:   req.Header.Filename,
		StoredPath: result.StoredPath,
		URL:        result.URL,
		MimeType:   saveReq.MimeType,
		SizeBytes:  req.Header.Size,
	}

	savedImg, err := s.uploadRepo.Save(ctx, coverImg)
	if err != nil {
		logger.Error("Failed to save cover to database",
			"series_id", req.SeriesID,
			"error", err.Error(),
		)

		// CRITICAL: Rollback - delete file from storage
		_ = s.storage.Delete(ctx, result.StoredPath)

		return param.UploadCoverResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to save cover record").
			WithKind(richerror.KindUnexpected)
	}

	logger.Info("Cover record saved to database",
		"series_id", req.SeriesID,
		"image_id", savedImg.ID,
	)

	// 5. Update series cover_image_url
	series.CoverImageURL = savedImg.URL
	_, err = s.seriesRepo.Update(ctx, req.SeriesID, series)
	if err != nil {
		logger.Error("Failed to update series cover URL",
			"series_id", req.SeriesID,
			"error", err.Error(),
		)

		// CRITICAL: Rollback - delete both DB record and file
		_ = s.uploadRepo.DeleteByID(ctx, savedImg.ID)
		_ = s.storage.Delete(ctx, result.StoredPath)

		return param.UploadCoverResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to update series").
			WithKind(richerror.KindUnexpected)
	}

	logger.Info("Cover uploaded successfully",
		"series_id", req.SeriesID,
		"cover_url", savedImg.URL,
	)

	return param.UploadCoverResponse{
		CoverImageURL: savedImg.URL,
	}, nil
}