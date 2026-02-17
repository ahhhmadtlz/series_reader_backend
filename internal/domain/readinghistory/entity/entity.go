package entity

import "time"

type ReadingHistory struct {
	ID       uint
	UserID   uint
	ChapterID  uint
	ReadAt   time.Time
} 