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

func (s Service) UploadBanner(ctx context.Context, req param.UploadBannerRequest) (param.UploadBannerResponse, error) {
	const op = richerror.Op("service.upload.UploadBanner")

	logger.Info("Upload banner request",
		"series_id", req.SeriesID,
		"user_id", req.UserID,
		"filename", req.Header.Filename,
		"size", req.Header.Size,
	)

	// 1. Check series exists
	_, err := s.seriesRepo.GetByID(ctx, req.SeriesID)
	if err != nil {
		return param.UploadBannerResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("series not found").
			WithKind(richerror.KindNotFound)
	}

	// 2. Delete old banner if exists (best effort)
	oldBanner, err := s.uploadRepo.GetLatestByOwner(ctx, req.SeriesID, entity.ImageKindBanner)
	if err == nil {
		_ = s.uploadRepo.DeleteByID(ctx, oldBanner.ID)
		if err := s.storage.Delete(ctx, oldBanner.StoredPath); err != nil {
			logger.Error("failed to delete old banner file",
				"series_id", req.SeriesID,
				"stored_path", oldBanner.StoredPath,
				"error", err,
			)
		}
	}

	// 3. Save to storage
	result, err := s.storage.Save(ctx, storage.SaveRequest{
		File:     req.File,
		Filename: req.Header.Filename,
		OwnerID:  req.SeriesID,
		Kind:     entity.ImageKindBanner,
		MimeType: req.Header.Header.Get("Content-Type"),
		Size:     req.Header.Size,
	})
	if err != nil {
		return param.UploadBannerResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to save banner file").
			WithKind(richerror.KindUnexpected)
	}

	// 4. Save record to DB
	bannerImg := entity.UploadedImage{
		OwnerID:    req.SeriesID,
		Kind:       entity.ImageKindBanner,
		Filename:   req.Header.Filename,
		StoredPath: result.StoredPath,
		URL:        result.URL,
		MimeType:   req.Header.Header.Get("Content-Type"),
		SizeBytes:  req.Header.Size,
	}

	savedImg, err := s.uploadRepo.Save(ctx, bannerImg)
	if err != nil {
		_ = s.storage.Delete(ctx, result.StoredPath)
		return param.UploadBannerResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to save banner record").
			WithKind(richerror.KindUnexpected)
	}

	logger.Info("Banner uploaded successfully",
		"series_id", req.SeriesID,
		"banner_url", savedImg.URL,
	)

	// 5. Enqueue image processing job (fire and forget)
	if err := s.jobQueue.Enqueue(ctx, ipparam.ProcessImageArgs{
		OwnerID:    req.SeriesID,
		RemotePath: savedImg.StoredPath,
		ImageKind:  entity.ImageKindBanner,
	}); err != nil {
		logger.Error("failed to enqueue banner processing job",
			"series_id", req.SeriesID,
			"error", err,
		)
	}

	return param.UploadBannerResponse{
		BannerImageURL: savedImg.URL,
	}, nil
}