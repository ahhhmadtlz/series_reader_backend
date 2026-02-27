package repository

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/entity"
	uploadentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/entity"
)

// Repository handles chapter page image variants (image_variants table).
type Repository interface {
	CreateVariant(ctx context.Context, variant entity.ImageVariant) (entity.ImageVariant, error)
	GetVariantsByPageID(ctx context.Context, pageID uint) ([]entity.ImageVariant, error)
	GetVariantsByPageIDs(ctx context.Context, pageIDs []uint) (map[uint][]entity.ImageVariant, error)
	DeleteVariantsByPageID(ctx context.Context, pageID uint) error
}

// VariantRepository is a unified interface for cover, banner, and chapter
// thumbnail variants. Each physical table keeps its own FK and ON DELETE
// CASCADE — this interface collapses the repeated method signatures into
// one contract using ImageKind + ownerID as discriminators.
//
// ownerID meaning per kind:
//
//	ImageKindCover            → series_id
//	ImageKindBanner           → series_id
//	ImageKindChapterThumbnail → chapter_id
type VariantRepository interface {
	CreateOwnerVariant(ctx context.Context, kind uploadentity.ImageKind, ownerID uint, variant entity.OwnerVariant) (entity.OwnerVariant, error)
	GetOwnerVariants(ctx context.Context, kind uploadentity.ImageKind, ownerID uint) ([]entity.OwnerVariant, error)
	DeleteOwnerVariants(ctx context.Context, kind uploadentity.ImageKind, ownerID uint) error
}

// CoverVariantRepository is kept for backwards compatibility during the
// transition to VariantRepository. New code should use VariantRepository.
//
// Deprecated: use VariantRepository with ImageKindCover.
type CoverVariantRepository interface {
	CreateCoverVariant(ctx context.Context, variant entity.CoverVariant) (entity.CoverVariant, error)
	GetCoverVariantsBySeriesID(ctx context.Context, seriesID uint) ([]entity.CoverVariant, error)
	DeleteCoverVariantsBySeriesID(ctx context.Context, seriesID uint) error
}

// BannerVariantRepository is kept for backwards compatibility during the
// transition to VariantRepository. New code should use VariantRepository.
//
// Deprecated: use VariantRepository with ImageKindBanner.
type BannerVariantRepository interface {
	CreateBannerVariant(ctx context.Context, variant entity.BannerVariant) (entity.BannerVariant, error)
	GetBannerVariantsBySeriesID(ctx context.Context, seriesID uint) ([]entity.BannerVariant, error)
	DeleteBannerVariantsBySeriesID(ctx context.Context, seriesID uint) error
}

// ChapterThumbnailRepository is kept for backwards compatibility during the
// transition to VariantRepository. New code should use VariantRepository.
//
// Deprecated: use VariantRepository with ImageKindChapterThumbnail.
type ChapterThumbnailRepository interface {
	CreateChapterThumbnailVariant(ctx context.Context, variant entity.ChapterThumbnailVariant) (entity.ChapterThumbnailVariant, error)
	GetChapterThumbnailVariantsByChapterID(ctx context.Context, chapterID uint) ([]entity.ChapterThumbnailVariant, error)
	DeleteChapterThumbnailVariantsByChapterID(ctx context.Context, chapterID uint) error
}