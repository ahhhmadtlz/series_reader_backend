package uploadhandler

import (
	"context"

	uploadparam "github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/param"
	uploadvalidator "github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/validator"
)

type UploadService interface {
    UploadAvatar(ctx context.Context, req uploadparam.UploadAvatarRequest) (uploadparam.UploadAvatarResponse, error)
    UploadCover(ctx context.Context, req uploadparam.UploadCoverRequest) (uploadparam.UploadCoverResponse, error)
    UploadBanner(ctx context.Context, req uploadparam.UploadBannerRequest) (uploadparam.UploadBannerResponse, error)
    UploadChapterThumbnail(ctx context.Context, req uploadparam.UploadChapterThumbnailRequest) (uploadparam.UploadChapterThumbnailResponse, error)
}

type Handler struct {
    service   UploadService
    validator uploadvalidator.Validator
}

func New(service UploadService, validator uploadvalidator.Validator) Handler {
    return Handler{
        service:   service,
        validator: validator,
    }
}