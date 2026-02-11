package service

import (
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/param"
)
func toChapterResponse(chapter *entity.Chapter) param.ChapterResponse {
	return param.ChapterResponse{
		ID:            chapter.ID,
		SeriesID:      chapter.SeriesID,
		ChapterNumber: chapter.ChapterNumber,
		Title:         chapter.Title,
	}
}

func toChapterWithPagesResponse(chapter *entity.Chapter, pages []entity.ChapterPage) param.ChapterWithPagesResponse {
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