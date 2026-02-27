package service

import (
	"context"

	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

// pathFetcher fetches the remote storage paths for all variants of an owner.
type pathFetcher func(ctx context.Context, id uint) ([]string, error)

// rowDeleter deletes all variant rows for an owner from the DB.
type rowDeleter func(ctx context.Context, id uint) error

// deleteVariants is the single algorithm for all variant deletion:
//  1. Fetch remote paths via fetch
//  2. Delete files from storage (failures are logged, never fatal)
//  3. Delete DB rows via delete
func (s Service) deleteVariants(
	ctx context.Context,
	op richerror.Op,
	id uint,
	fetch pathFetcher,
	delete rowDeleter,
	fetchErrMsg string,
	deleteErrMsg string,
) error {
	paths, err := fetch(ctx, id)
	if err != nil {
		return richerror.New(op).WithErr(err).
			WithMessage(fetchErrMsg).
			WithKind(richerror.KindUnexpected)
	}

	s.deleteVariantFiles(ctx, paths)

	if err := delete(ctx, id); err != nil {
		return richerror.New(op).WithErr(err).
			WithMessage(deleteErrMsg).
			WithKind(richerror.KindUnexpected)
	}

	return nil
}

// deleteVariantFiles removes files from storage for each remote path.
// Failures are logged and skipped — a missing file must never block
// the DB row deletion that follows.
//
// NOTE: files are deleted before DB rows. If file deletion partially
// fails and DB deletion succeeds, orphaned files may remain in storage.
// This is intentional: a broken DB row is worse than a leaked file.
// If consistency becomes critical, introduce a "pending_delete" marker
// and retry file deletion via a background worker.
func (s Service) deleteVariantFiles(ctx context.Context, paths []string) {
	for _, path := range paths {
		if path == "" {
			continue
		}
		if err := s.storage.Delete(ctx, path); err != nil {
			logger.Warn("failed to delete variant file, skipping", "path", path, "error", err)
		}
	}
}

func (s Service) DeletePageVariants(ctx context.Context, pageID uint) error {
	return s.deleteVariants(
		ctx,
		richerror.Op("service.imageprocessing.DeletePageVariants"),
		pageID,
		func(ctx context.Context, id uint) ([]string, error) {
			variants, err := s.repo.GetVariantsByPageID(ctx, id)
			if err != nil {
				return nil, err
			}
			paths := make([]string, 0, len(variants))
			for _, v := range variants {
				if v.RemotePath != "" {
					paths = append(paths, v.RemotePath)
				}
			}
			return paths, nil
		},
		s.repo.DeleteVariantsByPageID,
		"failed to fetch page variants",
		"failed to delete page variant rows",
	)
}

func (s Service) DeleteCoverVariants(ctx context.Context, seriesID uint) error {
	return s.deleteVariants(
		ctx,
		richerror.Op("service.imageprocessing.DeleteCoverVariants"),
		seriesID,
		func(ctx context.Context, id uint) ([]string, error) {
			variants, err := s.coverRepo.GetCoverVariantsBySeriesID(ctx, id)
			if err != nil {
				return nil, err
			}
			paths := make([]string, 0, len(variants))
			for _, v := range variants {
				if v.RemotePath != "" {
					paths = append(paths, v.RemotePath)
				}
			}
			return paths, nil
		},
		s.coverRepo.DeleteCoverVariantsBySeriesID,
		"failed to fetch cover variants",
		"failed to delete cover variant rows",
	)
}

func (s Service) DeleteBannerVariants(ctx context.Context, seriesID uint) error {
	return s.deleteVariants(
		ctx,
		richerror.Op("service.imageprocessing.DeleteBannerVariants"),
		seriesID,
		func(ctx context.Context, id uint) ([]string, error) {
			variants, err := s.bannerRepo.GetBannerVariantsBySeriesID(ctx, id)
			if err != nil {
				return nil, err
			}
			paths := make([]string, 0, len(variants))
			for _, v := range variants {
				if v.RemotePath != "" {
					paths = append(paths, v.RemotePath)
				}
			}
			return paths, nil
		},
		s.bannerRepo.DeleteBannerVariantsBySeriesID,
		"failed to fetch banner variants",
		"failed to delete banner variant rows",
	)
}

func (s Service) DeleteThumbnailVariants(ctx context.Context, chapterID uint) error {
	return s.deleteVariants(
		ctx,
		richerror.Op("service.imageprocessing.DeleteThumbnailVariants"),
		chapterID,
		func(ctx context.Context, id uint) ([]string, error) {
			variants, err := s.thumbnailRepo.GetChapterThumbnailVariantsByChapterID(ctx, id)
			if err != nil {
				return nil, err
			}
			paths := make([]string, 0, len(variants))
			for _, v := range variants {
				if v.RemotePath != "" {
					paths = append(paths, v.RemotePath)
				}
			}
			return paths, nil
		},
		s.thumbnailRepo.DeleteChapterThumbnailVariantsByChapterID,
		"failed to fetch thumbnail variants",
		"failed to delete thumbnail variant rows",
	)
}