package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) Delete(ctx context.Context, id uint) error {
	const op = richerror.Op("service.series.Delete")

	// 1. delete cover variant files
	coverVariants, err := s.coverRepo.GetCoverVariantsBySeriesID(ctx, id)
	if err != nil {
		logger.Error("failed to get cover variants", "series_id", id, "error", err)
	}
	for _, v := range coverVariants {
		if v.RemotePath != "" {
			if err := s.storage.Delete(ctx, v.RemotePath); err != nil {
				logger.Error("failed to delete cover variant file", "remote_path", v.RemotePath, "error", err)
			}
		}
	}

	// 2. delete banner variant files
	bannerVariants, err := s.bannerRepo.GetBannerVariantsBySeriesID(ctx, id)
	if err != nil {
		logger.Error("failed to get banner variants", "series_id", id, "error", err)
	}
	for _, v := range bannerVariants {
		if v.RemotePath != "" {
			if err := s.storage.Delete(ctx, v.RemotePath); err != nil {
				logger.Error("failed to delete banner variant file", "remote_path", v.RemotePath, "error", err)
			}
		}
	}

	// 3. delete cover and banner source images
	for _, kind := range []entity.ImageKind{entity.ImageKindCover, entity.ImageKindBanner} {
		images, err := s.uploadRepo.GetByOwner(ctx, id, kind)
		if err != nil {
			logger.Error("failed to get uploaded images", "series_id", id, "kind", kind, "error", err)
			continue
		}
		for _, img := range images {
			if err := s.storage.Delete(ctx, img.StoredPath); err != nil {
				logger.Error("failed to delete uploaded image file", "stored_path", img.StoredPath, "error", err)
			}
		}
	}

	// 4. delete all chapter pages, their variant files, thumbnail variant files, and source files
	chapters, err := s.chapterRepo.GetBySeriesID(ctx, id)
	if err != nil {
		logger.Error("failed to get chapters for series", "series_id", id, "error", err)
	}
	for _, ch := range chapters {
		// thumbnail variant files
		thumbnailVariants, err := s.thumbnailRepo.GetChapterThumbnailVariantsByChapterID(ctx, ch.ID)
		if err != nil {
			logger.Error("failed to get thumbnail variants", "chapter_id", ch.ID, "error", err)
		}
		for _, v := range thumbnailVariants {
			if v.RemotePath != "" {
				if err := s.storage.Delete(ctx, v.RemotePath); err != nil {
					logger.Error("failed to delete thumbnail variant file", "remote_path", v.RemotePath, "error", err)
				}
			}
		}

		// chapter thumbnail source image
		thumbImages, err := s.uploadRepo.GetByOwner(ctx, ch.ID, entity.ImageKindChapterThumbnail)
		if err != nil {
			logger.Error("failed to get thumbnail images", "chapter_id", ch.ID, "error", err)
		}
		for _, img := range thumbImages {
			if err := s.storage.Delete(ctx, img.StoredPath); err != nil {
				logger.Error("failed to delete thumbnail image file", "stored_path", img.StoredPath, "error", err)
			}
		}

		// page variant files and source files
		pages, err := s.chapterRepo.GetPagesByChapterID(ctx, ch.ID)
		if err != nil {
			logger.Error("failed to get pages", "chapter_id", ch.ID, "error", err)
			continue
		}
		for _, p := range pages {
			variants, err := s.imageVariantRepo.GetVariantsByPageID(ctx, p.ID)
			if err != nil {
				logger.Error("failed to get page variants", "page_id", p.ID, "error", err)
			}
			for _, v := range variants {
				if v.RemotePath != "" {
					if err := s.storage.Delete(ctx, v.RemotePath); err != nil {
						logger.Error("failed to delete page variant file", "remote_path", v.RemotePath, "error", err)
					}
				}
			}
			if p.RemotePath != "" {
				if err := s.storage.Delete(ctx, p.RemotePath); err != nil {
					logger.Error("failed to delete page source file", "remote_path", p.RemotePath, "error", err)
				}
			}
		}
	}

	// 5. delete the series — DB cascades handle all rows
	if err := s.repo.Delete(ctx, id); err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to delete series").
			WithKind(richerror.KindNotFound)
	}

	logger.Info("series deleted", "series_id", id)

	return nil
}