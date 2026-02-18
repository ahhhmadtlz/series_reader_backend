package entity

import (
	"time"

	sharedentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/shared/entity"
)

// ============================================
// Series Status & Type Constants
// ============================================

const (
	StatusOngoing   = "ongoing"
	StatusCompleted = "completed"
	StatusHiatus    = "hiatus"
	StatusCancelled = "cancelled"
)

const (
	TypeManga   = "manga"
	TypeManhwa  = "manhwa"
	TypeManhua  = "manhua"
	TypeComic   = "comic"
	TypeWebtoon = "webtoon"
)

// ============================================
// Series Entity
// ============================================

type Series struct {
	ID                uint
	Title             string
	Slug              string
	SlugID            string
	FullSlug          string
	Description       string
	Author            string
	Artist            string
	Status            string
	Type              string
	Genres            []string
	AlternativeTitles []string
	CoverImageURL     string
	PublicationYear   *int
	ViewCount         int
	Rating            float64
	IsPublished       bool
	IsPremiumOnly     bool  
	CreatedBy         *uint
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// ============================================
// SeriesCollaborator Entity
// ============================================

type SeriesCollaborator struct {
	SeriesID          uint
	UserID            uint
	CollaborationRole sharedentity.CollaborationRole
	AddedBy           uint
	AddedAt           time.Time
}

// ============================================
// Validation Helpers
// ============================================

func ValidStatuses() []string {
	return []string{StatusOngoing, StatusCompleted, StatusHiatus, StatusCancelled}
}

func ValidTypes() []string {
	return []string{TypeManga, TypeManhwa, TypeManhua, TypeComic, TypeWebtoon}
}

func IsValidStatus(status string) bool {
	for _, s := range ValidStatuses() {
		if s == status {
			return true
		}
	}
	return false
}

func IsValidType(t string) bool {
	for _, typ := range ValidTypes() {
		if typ == t {
			return true
		}
	}
	return false
}