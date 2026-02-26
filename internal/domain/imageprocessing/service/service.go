package service

import (
	iprepo "github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/repository"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/storage"
)

type Service struct {
	repo          iprepo.Repository
	coverRepo     iprepo.CoverVariantRepository
	bannerRepo    iprepo.BannerVariantRepository
	thumbnailRepo iprepo.ChapterThumbnailRepository
	storage       storage.Storage
}

func New(
	repo iprepo.Repository,
	coverRepo iprepo.CoverVariantRepository,
	bannerRepo iprepo.BannerVariantRepository,
	thumbnailRepo iprepo.ChapterThumbnailRepository,
	store storage.Storage,
) Service {
	return Service{
		repo:          repo,
		coverRepo:     coverRepo,
		bannerRepo:    bannerRepo,
		thumbnailRepo: thumbnailRepo,
		storage: store,
	}
}