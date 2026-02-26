package worker

import (
	"context"
	"os"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/param"
	iprepo "github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/repository"
	ipservice "github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/service"
	uploadentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/imageprocessor"
	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	"github.com/riverqueue/river"
)

type ImageProcessingWorker struct {
	river.WorkerDefaults[param.ProcessImageArgs]
	processor imageprocessor.Processor
	repo      iprepo.Repository
	coverRepo iprepo.CoverVariantRepository
	bannerRepo iprepo.BannerVariantRepository
	thumbnailRepo iprepo.ChapterThumbnailRepository
	ipSvc     ipservice.Service
	basePath  string
}

func NewImageProcessingWorker(
	processor imageprocessor.Processor,
	repo iprepo.Repository,
	coverRepo iprepo.CoverVariantRepository,
	bannerRepo iprepo.BannerVariantRepository,
	thumbnailRepo iprepo.ChapterThumbnailRepository,
	ipSvc ipservice.Service,
	basePath string,
) *ImageProcessingWorker {
	return &ImageProcessingWorker{
		processor:     processor,
		repo:          repo,
		coverRepo:     coverRepo,
		bannerRepo:    bannerRepo,
		thumbnailRepo: thumbnailRepo,
		ipSvc:         ipSvc,
		basePath:      basePath,
	}
}

func (w *ImageProcessingWorker) Work(ctx context.Context, job *river.Job[param.ProcessImageArgs]) error {
	const op = richerror.Op("worker.ImageProcessingWorker.Work")
	args := job.Args

	logger.Info("processing image variants", "kind", args.ImageKind, "owner_id", args.OwnerID)

	fullPath := w.basePath + "/" + args.RemotePath
	f, err := os.Open(fullPath)
	if err != nil {
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to open source image").
			WithKind(richerror.KindUnexpected)
	}
	defer f.Close()

	if err := w.deleteExistingVariants(ctx, args); err != nil {
		return err
	}

	variants := []entity.VariantKind{
		entity.VariantKindWebP,
		entity.VariantKindOptimized,
		entity.VariantKindCDN,
		entity.VariantKindThumbnail,
	}

	for _, variantKind := range variants {
		f.Seek(0, 0)

		result, err := w.processor.Process(ctx, f, args.PageID, args.OwnerID, args.ImageKind, variantKind)
		if err != nil {
			logger.Error("failed to process variant", "kind", variantKind, "error", err)
			return err
		}

		if err := w.saveVariant(ctx, args, result); err != nil {
			logger.Error("failed to save variant", "kind", variantKind, "error", err)
			return err
		}

		logger.Debug("variant saved", "kind", variantKind, "owner_id", args.OwnerID)
	}

	logger.Info("all variants processed", "owner_id", args.OwnerID)
	return nil
}

func (w *ImageProcessingWorker) saveVariant(ctx context.Context, args param.ProcessImageArgs, result param.ProcessImageResult) error {
	const op = richerror.Op("worker.ImageProcessingWorker.saveVariant")

	switch args.ImageKind {
	case uploadentity.ImageKindChapterPage:
		_, err := w.repo.CreateVariant(ctx, entity.ImageVariant{
			ChapterPageID: args.PageID,
			Kind:          result.Kind,
			ImageURL:      result.ImageURL,
			RemotePath:    result.RemotePath,
		})
		return err

	case uploadentity.ImageKindCover:
		_, err := w.coverRepo.CreateCoverVariant(ctx, entity.CoverVariant{
			SeriesID:   args.OwnerID,
			Kind:       result.Kind,
			ImageURL:   result.ImageURL,
			RemotePath: result.RemotePath,
		})
		return err

	case uploadentity.ImageKindBanner:
		_, err := w.bannerRepo.CreateBannerVariant(ctx, entity.BannerVariant{
			SeriesID:   args.OwnerID,
			Kind:       result.Kind,
			ImageURL:   result.ImageURL,
			RemotePath: result.RemotePath,
		})
		return err

	case uploadentity.ImageKindChapterThumbnail:
		_, err := w.thumbnailRepo.CreateChapterThumbnailVariant(ctx, entity.ChapterThumbnailVariant{
			ChapterID:  args.OwnerID,
			Kind:       result.Kind,
			ImageURL:   result.ImageURL,
			RemotePath: result.RemotePath,
		})
		return err

	default:
		return richerror.New(op).
			WithMessage("unknown image kind for variant saving").
			WithKind(richerror.KindInvalid)
	}
}

func (w *ImageProcessingWorker) deleteExistingVariants(ctx context.Context, args param.ProcessImageArgs) error {
	const op = richerror.Op("worker.ImageProcessingWorker.deleteExistingVariants")

	switch args.ImageKind {
	case uploadentity.ImageKindChapterPage:
		if err := w.ipSvc.DeletePageVariants(ctx, args.PageID); err != nil {
			return richerror.New(op).WithErr(err).WithMessage("failed to delete existing image variants").WithKind(richerror.KindUnexpected)
		}
	case uploadentity.ImageKindCover:
		if err := w.ipSvc.DeleteCoverVariants(ctx, args.OwnerID); err != nil {
			return richerror.New(op).WithErr(err).WithMessage("failed to delete existing cover variants").WithKind(richerror.KindUnexpected)
		}
	case uploadentity.ImageKindBanner:
		if err := w.ipSvc.DeleteBannerVariants(ctx, args.OwnerID); err != nil {
			return richerror.New(op).WithErr(err).WithMessage("failed to delete existing banner variants").WithKind(richerror.KindUnexpected)
		}
	case uploadentity.ImageKindChapterThumbnail:
		if err := w.ipSvc.DeleteThumbnailVariants(ctx, args.OwnerID); err != nil {
			return richerror.New(op).WithErr(err).WithMessage("failed to delete existing thumbnail variants").WithKind(richerror.KindUnexpected)
		}
	}

	return nil
}