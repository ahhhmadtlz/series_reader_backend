package entity

import "time"

type Chapter struct {
	ID            uint
	SeriesID      uint
	ChapterNumber float64
	Title        *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Pages []ChapterPage
}

type ChapterPage struct {
	ID uint
	ChapterID uint
	PageNumber int
	ImageURL string
	CreatedAt time.Time
}