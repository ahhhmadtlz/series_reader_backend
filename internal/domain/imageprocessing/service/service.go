package service

import iprepo "github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/repository"

// TODO: Add DeleteVariants methods (page/cover/banner/thumbnail) to this service.
// They should accept storage.Storage, fetch existing variant paths, delete files first,
// then delete DB rows. All callers (worker, chapter service, series service) should
// use this service instead of calling repos and storage directly.

type Service struct {
	repo          iprepo.Repository
	coverRepo     iprepo.CoverVariantRepository
	bannerRepo    iprepo.BannerVariantRepository
	thumbnailRepo iprepo.ChapterThumbnailRepository
}

func New(
	repo iprepo.Repository,
	coverRepo iprepo.CoverVariantRepository,
	bannerRepo iprepo.BannerVariantRepository,
	thumbnailRepo iprepo.ChapterThumbnailRepository,
) Service {
	return Service{
		repo:          repo,
		coverRepo:     coverRepo,
		bannerRepo:    bannerRepo,
		thumbnailRepo: thumbnailRepo,
	}
}