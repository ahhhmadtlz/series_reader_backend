package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) GetByID(ctx context.Context, id uint) (param.ChapterResponse, error) {
	const op = richerror.Op("service.chapter.GetByID")

	chapter, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return param.ChapterResponse{}, richerror.New(op).WithErr(err).WithMessage("failed to get chapter").WithKind(richerror.KindNotFound)
	}

	return toChapterResponse(chapter), nil
}

func (s Service) GetBySeriesID(ctx context.Context, seriesID uint) ([]param.ChapterResponse, error) {
	const op = richerror.Op("service.chapter.GetSeriesID")

	chapters, err := s.repo.GetBySeriesID(ctx, seriesID)
	if err != nil {
		return nil, richerror.New(op).WithErr(err).WithMessage("failed to get chapters for series").WithKind(richerror.KindUnexpected)
	}

	responses := make([]param.ChapterResponse, len(chapters))
	for i, ch := range chapters {
		responses[i] = toChapterResponse(ch)
	}
	return responses, nil
}

func (s Service) GetChapterWithPages(ctx context.Context, chapterID uint) (param.ChapterWithPagesResponse, error) {
	const op = richerror.Op("service.chapter.GetChapterWithPages")

	chapter, pages, err := s.repo.GetChapterWithPages(ctx, chapterID)
	if err != nil {
		return param.ChapterWithPagesResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get chapter with pages").
			WithKind(richerror.KindNotFound)
	}

	pageResponses := make([]param.ChapterPageResponse, len(pages))
	for i, p := range pages {
		variants, err := s.variantRepo.GetVariantsByPageID(ctx, p.ID)
		if err != nil {
			// fail silently — variants may not exist yet for freshly uploaded pages
			logger.Error("failed to get variants for page",
				"page_id", p.ID,
				"error", err,
			)
			variants = nil
		}

		variantResponses := make([]param.PageVariantResponse, len(variants))
		for j, v := range variants {
			variantResponses[j] = param.PageVariantResponse{
				ID:        v.ID,
				Kind:      v.Kind,
				ImageURL:  v.ImageURL,
				CreatedAt: v.CreatedAt.Format("2006-01-02T15:04:05Z"),
			}
		}

		pageResponses[i] = param.ChapterPageResponse{
			ID:         p.ID,
			PageNumber: p.PageNumber,
			ImageURL:   p.ImageURL,
			Variants:   variantResponses,
		}
	}

	return param.ChapterWithPagesResponse{
		ID:            chapter.ID,
		SeriesID:      chapter.SeriesID,
		ChapterNumber: chapter.ChapterNumber,
		Title:         chapter.Title,
		Pages:         pageResponses,
	}, nil
}