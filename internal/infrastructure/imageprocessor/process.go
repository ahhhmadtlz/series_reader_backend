package imageprocessor

import (
	"bytes"
	"context"
	"io"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/param"
	uploadentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

func (p *ImageProcessor) Process(ctx context.Context, src io.Reader, pageID uint, ownerID uint, imageKind uploadentity.ImageKind, variantKind entity.VariantKind) (param.ProcessImageResult, error) {
	const op = richerror.Op("imageprocessor.Process")

	data, err := io.ReadAll(src)
	if err != nil {
		return param.ProcessImageResult{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to read source image").
			WithKind(richerror.KindUnexpected)
	}

	switch variantKind {
	case entity.VariantKindWebP:
		return p.processWebP(ctx, bytes.NewReader(data), pageID, ownerID, imageKind)
	case entity.VariantKindThumbnail:
		return p.processThumbnail(ctx, bytes.NewReader(data), pageID, ownerID, imageKind)
	case entity.VariantKindOptimized:
		return p.processOptimized(ctx, bytes.NewReader(data), pageID, ownerID, imageKind)
	case entity.VariantKindCDN:
		return p.processCDN(ctx, bytes.NewReader(data), pageID, ownerID, imageKind)
	default:
		return param.ProcessImageResult{}, richerror.New(op).
			WithMessage("unknown variant kind").
			WithKind(richerror.KindInvalid)
	}
}