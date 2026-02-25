package param

import (
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/entity"
	uploadentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/entity"
)

// ProcessImageArgs are the job arguments passed to the River worker
type ProcessImageArgs struct {
	PageID     uint
	OwnerID    uint
	RemotePath string
	ChapterID  uint
	ImageKind  uploadentity.ImageKind
}

func (ProcessImageArgs) Kind() string { return "process_image" }

// ProcessImageResult is what the processor returns per variant
type ProcessImageResult struct {
	Kind       entity.VariantKind
	ImageURL   string
	RemotePath string
}

// ImageVariantResponse is the API response for a single image variant
type ImageVariantResponse struct {
	ID            uint               `json:"id"`
	ChapterPageID uint               `json:"chapter_page_id"`
	Kind          entity.VariantKind `json:"kind"`
	ImageURL      string             `json:"image_url"`
	CreatedAt     string             `json:"created_at"`
}

// GetVariantsResponse is the API response for listing all variants of a page
type GetVariantsResponse struct {
	PageID   uint                   `json:"page_id"`
	Variants []ImageVariantResponse `json:"variants"`
}


// CoverVariantResponse is the API response for a single cover variant
type CoverVariantResponse struct {
	ID         uint               `json:"id"`
	SeriesID   uint               `json:"series_id"`
	Kind       entity.VariantKind `json:"kind"`
	ImageURL   string             `json:"image_url"`
	CreatedAt  string             `json:"created_at"`
}

// GetCoverVariantsResponse is the API response for listing all cover variants of a series
type GetCoverVariantsResponse struct {
	SeriesID uint                   `json:"series_id"`
	Variants []CoverVariantResponse `json:"variants"`
}

// BannerVariantResponse is the API response for a single banner variant
type BannerVariantResponse struct {
	ID        uint               `json:"id"`
	SeriesID  uint               `json:"series_id"`
	Kind      entity.VariantKind `json:"kind"`
	ImageURL  string             `json:"image_url"`
	CreatedAt string             `json:"created_at"`
}

// GetBannerVariantsResponse is the API response for listing all banner variants of a series
type GetBannerVariantsResponse struct {
	SeriesID uint                    `json:"series_id"`
	Variants []BannerVariantResponse `json:"variants"`
}

// ThumbnailVariantResponse is the API response for a single chapter thumbnail variant
type ThumbnailVariantResponse struct {
	ID        uint               `json:"id"`
	ChapterID uint               `json:"chapter_id"`
	Kind      entity.VariantKind `json:"kind"`
	ImageURL  string             `json:"image_url"`
	CreatedAt string             `json:"created_at"`
}

// GetThumbnailVariantsResponse is the API response for listing all thumbnail variants of a chapter
type GetThumbnailVariantsResponse struct {
	ChapterID uint                       `json:"chapter_id"`
	Variants  []ThumbnailVariantResponse `json:"variants"`
}