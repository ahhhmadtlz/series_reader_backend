package param

import "time"

type MarkAsReadRequest struct {
	ChapterID uint `json:"chapter_id"`
}

type ReadingHistoryResponse struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	ChapterID     uint      `json:"chapter_id"`
	SeriesID      uint      `json:"series_id"`
  ChapterNumber float64 `json:"chapter_number"`
	SeriesTitle   string    `json:"series_title"`
	ReadAt        time.Time `json:"read_at"`
}


type SeriesProgressResponse struct {
	SeriesID       uint                  `json:"series_id"`
	SeriesTitle    string                `json:"series_title"`
	TotalChapters  int                   `json:"total_chapters"`
	ReadChapters   int                   `json:"read_chapters"`
	Chapters       []ChapterProgressItem `json:"chapters"`
}


type ChapterProgressItem struct {
	ChapterID     uint       `json:"chapter_id"`
	ChapterNumber float64    `json:"chapter_number"`
	IsRead        bool       `json:"is_read"`
	ReadAt        *time.Time `json:"read_at,omitempty"`
}