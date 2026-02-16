package entity

import "time"

type Bookmark struct {
	ID        uint
	UserID    uint
	SeriesID    uint
	CreatedAt time.Time
}
