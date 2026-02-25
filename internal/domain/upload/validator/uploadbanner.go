package validator

import (
	"context"
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateUploadBanner(ctx context.Context, req param.UploadBannerRequest) error {
	const op = richerror.Op("validator.upload.ValidateUploadBanner")

	fieldErrors := make(map[string]string)

	err := validation.ValidateStruct(&req,
		validation.Field(&req.SeriesID,
			validation.Required.Error("series ID is required"),
		),
		validation.Field(&req.UserID,
			validation.Required.Error("user ID is required"),
		),
		validation.Field(&req.Header,
			validation.Required.Error("file is required"),
			validation.By(v.validateBannerFile),
		),
	)

	if err != nil {
		if errV, ok := err.(validation.Errors); ok {
			for field, value := range errV {
				if value != nil {
					fieldErrors[field] = value.Error()
				}
			}
		}
	}

	if len(fieldErrors) > 0 {
		return richerror.New(op).
			WithMessage("invalid input").
			WithKind(richerror.KindInvalid).
			WithMeta("fields", fieldErrors)
	}

	return nil
}

func (v Validator) validateBannerFile(value interface{}) error {
	header, ok := value.(*multipart.FileHeader)
	if !ok || header == nil {
		return errors.New("invalid file header")
	}

	if header.Size == 0 {
		return errors.New("file is empty")
	}

	maxSizeBytes := int64(v.uploadConfig.MaxBannerSizeMB) * 1024 * 1024
	if header.Size > maxSizeBytes {
		return errors.New("banner size exceeds limit")
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		return errors.New("Content-Type header is required")
	}

	if !isAllowedMimeType(mimeType, v.uploadConfig.AllowedMimeTypes) {
		return errors.New("invalid file type. Allowed: " + strings.Join(v.uploadConfig.AllowedMimeTypes, ", "))
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !isAllowedExtension(ext) {
		return errors.New("invalid file extension. Allowed: .jpg, .jpeg, .png, .webp")
	}

	if err := validateFilename(header.Filename); err != nil {
		return err
	}

	return nil
}