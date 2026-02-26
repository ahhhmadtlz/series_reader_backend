package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) DeletePageVariants(ctx context.Context, pageID uint) error{
	const op=richerror.Op("service.imageprocessing.DeletePageVariants")

	variants,err :=s.repo.GetVariantsByPageID(ctx,pageID)
	if err !=nil{
		return richerror.New(op).
		 WithErr(err).
		 WithMessage("failed to fetch page variants").
		 WithKind(richerror.KindUnexpected)
	}

  for _, v := range variants {
		if v.RemotePath != "" {
			if err := s.storage.Delete(ctx, v.RemotePath); err != nil {
				logger.Warn("failed to delete page variant file, skipping", "path", v.RemotePath, "error", err)
			}
		}
	}

	if err := s.repo.DeleteVariantsByPageID(ctx, pageID); err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to delete page variant rows").
			WithKind(richerror.KindUnexpected)
	}

	return nil
}

func (s Service) DeleteCoverVariants(ctx context.Context, seriesID uint) error {
	const op = richerror.Op("service.imageprocessing.DeleteCoverVariants")

	variants, err := s.coverRepo.GetCoverVariantsBySeriesID(ctx, seriesID)
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to fetch cover variants").
			WithKind(richerror.KindUnexpected)
	}

	for _, v := range variants {
		if v.RemotePath != "" {
			if err := s.storage.Delete(ctx, v.RemotePath); err != nil {
				logger.Warn("failed to delete cover variant file, skipping", "path", v.RemotePath, "error", err)
			}
		}
	}

	if err := s.coverRepo.DeleteCoverVariantsBySeriesID(ctx, seriesID); err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to delete cover variant rows").
			WithKind(richerror.KindUnexpected)
	}

	return nil
}

func (s Service) DeleteBannerVariants(ctx context.Context, seriesID uint) error {
	const op = richerror.Op("service.imageprocessing.DeleteBannerVariants")

	variants, err := s.bannerRepo.GetBannerVariantsBySeriesID(ctx, seriesID)
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to fetch banner variants").
			WithKind(richerror.KindUnexpected)
	}

	for _, v := range variants {
		if v.RemotePath != "" {
			if err := s.storage.Delete(ctx, v.RemotePath); err != nil {
				logger.Warn("failed to delete banner variant file, skipping", "path", v.RemotePath, "error", err)
			}
		}
	}

	if err := s.bannerRepo.DeleteBannerVariantsBySeriesID(ctx, seriesID); err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to delete banner variant rows").
			WithKind(richerror.KindUnexpected)
	}

	return nil
}


func (s Service) DeleteThumbnailVariants(ctx context.Context, chapterID uint) error {
	const op = richerror.Op("service.imageprocessing.DeleteThumbnailVariants")

	variants, err := s.thumbnailRepo.GetChapterThumbnailVariantsByChapterID(ctx, chapterID)
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to fetch thumbnail variants").
			WithKind(richerror.KindUnexpected)
	}

	for _, v := range variants {
		if v.RemotePath != "" {
			if err := s.storage.Delete(ctx, v.RemotePath); err != nil {
				logger.Warn("failed to delete thumbnail variant file, skipping", "path", v.RemotePath, "error", err)
			}
		}
	}

	if err := s.thumbnailRepo.DeleteChapterThumbnailVariantsByChapterID(ctx, chapterID); err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to delete thumbnail variant rows").
			WithKind(richerror.KindUnexpected)
	}

	return nil
}