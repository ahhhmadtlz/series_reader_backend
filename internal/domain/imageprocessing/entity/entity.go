package entity

import "time"

type VariantKind string

const (
	VariantKindWebP      VariantKind = "webp"
	VariantKindThumbnail VariantKind = "thumbnail"
	VariantKindOptimized VariantKind = "optimized"
	VariantKindCDN       VariantKind = "cdn"
)

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