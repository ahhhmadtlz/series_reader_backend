package imageprocessinghandler

import (
	"context"

	ipparam "github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/param"
)

type ImageProcessingService interface {
    GetVariants(ctx context.Context, pageID uint) (ipparam.GetVariantsResponse, error)
    GetCoverVariants(ctx context.Context, seriesID uint) (ipparam.GetCoverVariantsResponse, error)
    GetBannerVariants(ctx context.Context, seriesID uint) (ipparam.GetBannerVariantsResponse, error)
    GetThumbnailVariants(ctx context.Context, chapterID uint) (ipparam.GetThumbnailVariantsResponse, error)
}

type Handler struct {
    ipService ImageProcessingService
}

func New(ipService ImageProcessingService) Handler {
    return Handler{
        ipService: ipService,
    }
}