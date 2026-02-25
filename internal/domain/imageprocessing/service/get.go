package service

import (
	"context"
	"time"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) GetVariants(ctx context.Context, pageID uint) (param.GetVariantsResponse, error) {
	const op = richerror.Op("service.imageprocessing.GetVariants")

	variants, err := s.repo.GetVariantsByPageID(ctx, pageID)
	if err != nil {
		return param.GetVariantsResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get variants").
			WithKind(richerror.KindUnexpected)
	}

	responses := make([]param.ImageVariantResponse, len(variants))
	for i, v := range variants {
		responses[i] = param.ImageVariantResponse{
			ID:            v.ID,
			ChapterPageID: v.ChapterPageID,
			Kind:          v.Kind,
			ImageURL:      v.ImageURL,
			CreatedAt:     v.CreatedAt.UTC().Format(time.RFC3339),
		}
	}

	return param.GetVariantsResponse{
		PageID:   pageID,
		Variants: responses,
	}, nil
}


func (s Service) GetCoverVariants(ctx context.Context, seriesID uint) (param.GetCoverVariantsResponse, error) {
	const op = richerror.Op("service.imageprocessing.GetCoverVariants")

	variants, err := s.coverRepo.GetCoverVariantsBySeriesID(ctx, seriesID)
	if err != nil {
		return param.GetCoverVariantsResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get cover variants").
			WithKind(richerror.KindUnexpected)
	}

	responses := make([]param.CoverVariantResponse, len(variants))
	for i, v := range variants {
		responses[i] = param.CoverVariantResponse{
			ID:        v.ID,
			SeriesID:  v.SeriesID,
			Kind:      v.Kind,
			ImageURL:  v.ImageURL,
			CreatedAt:v.CreatedAt.UTC().Format(time.RFC3339),
		}
	}

	return param.GetCoverVariantsResponse{
		SeriesID: seriesID,
		Variants: responses,
	}, nil
}

func (s Service) GetBannerVariants(ctx context.Context, seriesID uint) (param.GetBannerVariantsResponse, error) {
	const op = richerror.Op("service.imageprocessing.GetBannerVariants")

	variants, err := s.bannerRepo.GetBannerVariantsBySeriesID(ctx, seriesID)
	if err != nil {
		return param.GetBannerVariantsResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get banner variants").
			WithKind(richerror.KindUnexpected)
	}

	responses := make([]param.BannerVariantResponse, len(variants))
	for i, v := range variants {
		responses[i] = param.BannerVariantResponse{
			ID:        v.ID,
			SeriesID:  v.SeriesID,
			Kind:      v.Kind,
			ImageURL:  v.ImageURL,
			CreatedAt: v.CreatedAt.UTC().Format(time.RFC3339),
		}
	}

	return param.GetBannerVariantsResponse{
		SeriesID: seriesID,
		Variants: responses,
	}, nil
}

func (s Service) GetThumbnailVariants(ctx context.Context, chapterID uint) (param.GetThumbnailVariantsResponse, error) {
	const op = richerror.Op("service.imageprocessing.GetThumbnailVariants")

	variants, err := s.thumbnailRepo.GetChapterThumbnailVariantsByChapterID(ctx, chapterID)
	if err != nil {
		return param.GetThumbnailVariantsResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get thumbnail variants").
			WithKind(richerror.KindUnexpected)
	}

	responses := make([]param.ThumbnailVariantResponse, len(variants))
	for i, v := range variants {
		responses[i] = param.ThumbnailVariantResponse{
			ID:        v.ID,
			ChapterID: v.ChapterID,
			Kind:      v.Kind,
			ImageURL:  v.ImageURL,
			CreatedAt: v.CreatedAt.UTC().Format(time.RFC3339),
		}
	}

	return param.GetThumbnailVariantsResponse{
		ChapterID: chapterID,
		Variants:  responses,
	}, nil
}