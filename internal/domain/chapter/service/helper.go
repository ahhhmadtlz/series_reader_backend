package service

import (
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/param"
	ipentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/entity"
)

func toChapterResponse(chapter *entity.Chapter) param.ChapterResponse {
	if chapter == nil {
		return param.ChapterResponse{}
	}

	return param.ChapterResponse{
		ID:            chapter.ID,
		SeriesID:      chapter.SeriesID,
		ChapterNumber: chapter.ChapterNumber,
		Title:         chapter.Title,
	}
}

func toChapterWithPagesResponse(chapter *entity.Chapter, pages []entity.ChapterPage) param.ChapterWithPagesResponse {
	if chapter == nil {
		return param.ChapterWithPagesResponse{}
	}

	pageResponses := make([]param.ChapterPageResponse, len(pages))
	for i, p := range pages {
		pageResponses[i] = param.ChapterPageResponse{
			ID:         p.ID,
			PageNumber: p.PageNumber,
			ImageURL:   p.ImageURL,
		}
	}

	return param.ChapterWithPagesResponse{
		ID:            chapter.ID,
		SeriesID:      chapter.SeriesID,
		ChapterNumber: chapter.ChapterNumber,
		Title:         chapter.Title,
		Pages:         pageResponses,
	}
}

func toChapterPageResponse(page *entity.ChapterPage) param.ChapterPageResponse {
	if page == nil {
		return param.ChapterPageResponse{}
	}

	return param.ChapterPageResponse{
		ID:         page.ID,
		PageNumber: page.PageNumber,
		ImageURL:   page.ImageURL,
	}
}


// buildPageResponse maps a page and its pre-fetched variants into a response.
// variantsByPage may be nil (e.g. fresh upload with no variants yet) — handled gracefully.
func buildPageResponse(p entity.ChapterPage, variantsByPage map[uint][]ipentity.ImageVariant) param.ChapterPageResponse {
	variants := variantsByPage[p.ID]
	variantResponses := make([]param.PageVariantResponse, len(variants))
	for i, v := range variants {
		variantResponses[i] = param.PageVariantResponse{
			ID:        v.ID,
			Kind:      v.Kind,
			ImageURL:  v.ImageURL,
			CreatedAt: v.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return param.ChapterPageResponse{
		ID:         p.ID,
		PageNumber: p.PageNumber,
		ImageURL:   p.ImageURL,
		Variants:   variantResponses,
	}
}