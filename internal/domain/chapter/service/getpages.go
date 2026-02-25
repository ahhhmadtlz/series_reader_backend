package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) GetPages(ctx context.Context, chapterID uint) ([]param.ChapterPageResponse, error) {
	const op = richerror.Op("service.chapter.GetPages")

	pages, err := s.repo.GetPagesByChapterID(ctx, chapterID)
	if err != nil {
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get pages").
			WithKind(richerror.KindUnexpected)
	}

	responses := make([]param.ChapterPageResponse, len(pages))
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

		responses[i] = param.ChapterPageResponse{
			ID:         p.ID,
			PageNumber: p.PageNumber,
			ImageURL:   p.ImageURL,
			Variants:   variantResponses,
		}
	}

	return responses, nil
}

func (s Service) GetPageByNumber(ctx context.Context, chapterID uint, pageNumber int) (param.ChapterPageResponse, error) {
	const op = richerror.Op("service.chapter.GetPageByNumber")

	page, err := s.repo.GetPageByNumber(ctx, chapterID, pageNumber)
	if err != nil {
		return param.ChapterPageResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("page not found").
			WithKind(richerror.KindNotFound)
	}

	variants, err := s.variantRepo.GetVariantsByPageID(ctx, page.ID)
	if err != nil {
		// fail silently — variants may not exist yet for freshly uploaded pages
		logger.Error("failed to get variants for page",
			"page_id", page.ID,
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

	return param.ChapterPageResponse{
		ID:         page.ID,
		PageNumber: page.PageNumber,
		ImageURL:   page.ImageURL,
		Variants:   variantResponses,
	}, nil
}