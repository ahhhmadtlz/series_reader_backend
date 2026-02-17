package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/readinghistory/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (s Service) GetSeriesProgress(ctx context.Context, userID uint, seriesSlug string) (param.SeriesProgressResponse, error) {
	const op = richerror.Op("service.readinghistory.GetSeriesProgress")

 	series, err := s.seriesRepo.GetByFullSlug(ctx, seriesSlug)
	if err != nil {
		return param.SeriesProgressResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("series not found").
			WithKind(richerror.KindNotFound)
	}

	allChapters, err := s.chapterRepo.GetBySeriesID(ctx, series.ID)
	if err != nil {
		return param.SeriesProgressResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get series chapters")
	}

	readHistories, err := s.readinghistoryRepo.GetSeriesProgress(ctx, userID, series.ID)
	if err != nil {
		return param.SeriesProgressResponse{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to get series progress")
	}

	readMap := make(map[uint]param.ReadingHistoryResponse)
	for _, h := range readHistories {
		// Find the chapter to get chapter_number
		for _, ch := range allChapters {
			if ch.ID == h.ChapterID {
				readMap[h.ChapterID] = param.ReadingHistoryResponse{
					ID:            h.ID,
					UserID:        h.UserID,
					ChapterID:     h.ChapterID,
					SeriesID:      series.ID,
					ChapterNumber: ch.ChapterNumber,
					SeriesTitle:   series.Title,
					ReadAt:        h.ReadAt,
				}
				break
			}
		}
	}

	chapters := make([]param.ChapterProgressItem, 0, len(allChapters))
	for _, ch := range allChapters {
		item := param.ChapterProgressItem{
			ChapterID:     ch.ID,
			ChapterNumber: ch.ChapterNumber,
			IsRead:        false,
		}

		// Check if this chapter is read
		if readInfo, exists := readMap[ch.ID]; exists {
			item.IsRead = true
			item.ReadAt = &readInfo.ReadAt
		}

		chapters = append(chapters, item)
	}
	response := param.SeriesProgressResponse{
		SeriesID:      series.ID,
		SeriesTitle:   series.Title,
		TotalChapters: len(allChapters),
		ReadChapters:  len(readHistories),
		Chapters:      chapters,
	}

	return response, nil
}