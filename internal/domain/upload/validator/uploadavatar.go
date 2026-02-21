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

func (v Validator) ValidateUploadAvatar(ctx context.Context, req param.UploadAvatarRequest) error {
	const op = richerror.Op("validator.upload.ValidateUploadAvatar")

	fieldErrors := make(map[string]string)

	err := validation.ValidateStruct(&req,
		validation.Field(&req.UserID,
			validation.Required.Error("user ID is required"),
		),
		validation.Field(&req.Header,
			validation.Required.Error("file is required"),
			validation.By(v.validateAvatarFile),
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

// validateAvatarFile is a custom validation rule for avatar files
func (v Validator) validateAvatarFile(value interface{}) error {
	header, ok := value.(*multipart.FileHeader)
	if !ok || header == nil {
		return errors.New("invalid file header")
	}

	// Check file not empty
	if header.Size == 0 {
		return errors.New("file is empty")
	}

	// Check domain-specific size limit (avatar)
	maxSizeBytes := int64(v.uploadConfig.MaxAvatarSizeMB) * 1024 * 1024
	if header.Size > maxSizeBytes {
		return errors.New("avatar size exceeds limit")
	}

	// Validate MIME type is in allowed list
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		return errors.New("Content-Type header is required")
	}
	
	if !isAllowedMimeType(mimeType, v.uploadConfig.AllowedMimeTypes) {
		return errors.New("invalid file type. Allowed: " + strings.Join(v.uploadConfig.AllowedMimeTypes, ", "))
	}

	// Validate extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !isAllowedExtension(ext) {
		return errors.New("invalid file extension. Allowed: .jpg, .jpeg, .png, .webp")
	}

	// Validate filename security
	if err := validateFilename(header.Filename); err != nil {
		return err
	}

	return nil
}