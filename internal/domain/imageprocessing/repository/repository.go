package repository

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/entity"
)

type Repository interface {
	CreateVariant(ctx context.Context, variant entity.ImageVariant) (entity.ImageVariant, error)
	GetVariantsByPageID(ctx context.Context, pageID uint) ([]entity.ImageVariant, error)
	DeleteVariantsByPageID(ctx context.Context, pageID uint) error
}

type CoverVariantRepository interface {
	CreateCoverVariant(ctx context.Context, variant entity.CoverVariant) (entity.CoverVariant, error)
	GetCoverVariantsBySeriesID(ctx context.Context, seriesID uint) ([]entity.CoverVariant, error)
	DeleteCoverVariantsBySeriesID(ctx context.Context, seriesID uint) error
}

type BannerVariantRepository interface {
	CreateBannerVariant(ctx context.Context, variant entity.BannerVariant) (entity.BannerVariant, error)
	GetBannerVariantsBySeriesID(ctx context.Context, seriesID uint) ([]entity.BannerVariant, error)
	DeleteBannerVariantsBySeriesID(ctx context.Context, seriesID uint) error
}

type ChapterThumbnailRepository interface {
	CreateChapterThumbnailVariant(ctx context.Context, variant entity.ChapterThumbnailVariant) (entity.ChapterThumbnailVariant, error)
	GetChapterThumbnailVariantsByChapterID(ctx context.Context, chapterID uint) ([]entity.ChapterThumbnailVariant, error)
	DeleteChapterThumbnailVariantsByChapterID(ctx context.Context, chapterID uint) error
}