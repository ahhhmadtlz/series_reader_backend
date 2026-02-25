package imageprocessor

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/param"
	uploadentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/storage"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	"github.com/disintegration/imaging"
)

func (p *ImageProcessor) processWebP(ctx context.Context, src io.Reader, pageID uint, ownerID uint, imageKind uploadentity.ImageKind) (param.ProcessImageResult, error) {
	const op = richerror.Op("imageprocessor.processWebP")

	img, err := imaging.Decode(src)
	if err != nil {
		return param.ProcessImageResult{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to decode image").
			WithKind(richerror.KindUnexpected)
	}

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, img, imaging.JPEG, imaging.JPEGQuality(85)); err != nil {
		return param.ProcessImageResult{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to encode image").
			WithKind(richerror.KindUnexpected)
	}

	result, err := p.storage.Save(ctx, storage.SaveRequest{
		File:     &buf,
		Filename: fmt.Sprintf("conversions/%d-webp.jpg", pageID),
		OwnerID:  ownerID,
		Kind:     imageKind,
		MimeType: "image/jpeg",
		Size:     int64(buf.Len()),
	})
	if err != nil {
		return param.ProcessImageResult{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to save webp variant").
			WithKind(richerror.KindUnexpected)
	}

	return param.ProcessImageResult{
		Kind:       entity.VariantKindWebP,
		ImageURL:   result.URL,
		RemotePath: result.StoredPath,
	}, nil
}