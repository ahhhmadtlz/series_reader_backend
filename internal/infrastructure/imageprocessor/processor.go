package imageprocessor

import (
	"context"
	"io"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/param"
	uploadentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/entity"
)

type Processor interface {
	Process(ctx context.Context, src io.Reader, pageID uint, ownerID uint, imageKind uploadentity.ImageKind, variantKind entity.VariantKind) (param.ProcessImageResult, error)
}