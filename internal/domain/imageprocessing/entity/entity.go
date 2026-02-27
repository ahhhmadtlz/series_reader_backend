package entity

import "time"

type VariantKind string

const (
	VariantKindWebP      VariantKind = "webp"
	VariantKindThumbnail VariantKind = "thumbnail"
	VariantKindOptimized VariantKind = "optimized"
	VariantKindCDN       VariantKind = "cdn"
)

// OwnerVariant is the shared value type used by VariantRepository.
// It replaces CoverVariant, BannerVariant, and ChapterThumbnailVariant
// at the interface boundary — the physical tables and their FKs are unchanged.
type OwnerVariant struct {
	ID         uint
	OwnerID    uint
	Kind       VariantKind
	ImageURL   string
	RemotePath string
	CreatedAt  time.Time
}

type ImageVariant struct {
	ID            uint
	ChapterPageID uint
	Kind          VariantKind
	ImageURL      string
	RemotePath    string
	CreatedAt     time.Time
}

type CoverVariant struct {
	ID            uint
	SeriesID      uint
	Kind          VariantKind
	ImageURL      string
	RemotePath    string
	CreatedAt     time.Time
}

type BannerVariant struct {
	ID            uint
	SeriesID      uint
	Kind          VariantKind
	ImageURL      string
	RemotePath    string
	CreatedAt     time.Time
}


type ChapterThumbnailVariant  struct{
	ID            uint
	ChapterID     uint
	Kind          VariantKind
	ImageURL      string
	RemotePath    string
	CreatedAt     time.Time
}