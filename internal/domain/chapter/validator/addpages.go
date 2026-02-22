package validator
import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/chapter/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)


func (v Validator) ValidateUploadPage(ctx context.Context, req param.UploadPageParam) error {

	const op = richerror.Op("validator.chapter.ValidateUploadPage")

	fieldErrors := make(map[string]string)

	if req.ChapterID == 0 {
		fieldErrors["chapter_id"] = "chapter ID is required"
	}

	if req.PageNumber < 0 {
		fieldErrors["page_number"] = "page number must be non-negative"
	}

	if req.Header == nil {
		fieldErrors["file"] = "file is required"
	} else {
		if err := validatePageFile(req.Header, v.uploadConfig.MaxPageSizeMB, v.uploadConfig.AllowedMimeTypes); err != nil {
			fieldErrors["file"] = err.Error()
		}
	}

	if len(fieldErrors) > 0 {
		return richerror.New(op).
			WithMessage("validation failed").
			WithKind(richerror.KindInvalid).
			WithMetaMap(toAnyMap(fieldErrors))
	}

	return nil
}

func (v Validator) ValidateBulkUpload(ctx context.Context, req param.BulkUploadParam) error {
	const op = richerror.Op("validator.chapter.ValidateBulkUpload")

	fieldErrors := make(map[string]string)

	if req.ChapterID == 0 {
		fieldErrors["chapter_id"] = "chapter ID is required"
	}

	if len(req.Files) == 0 {
		fieldErrors["files"] = "at least one file is required"
	}

	for i, fh := range req.Files {
		if err := validatePageFile(fh, v.uploadConfig.MaxPageSizeMB, v.uploadConfig.AllowedMimeTypes); err != nil {
			fieldErrors[fmt.Sprintf("files[%d]", i)] = err.Error()
		}
	}

	if len(fieldErrors) > 0 {
		return richerror.New(op).
			WithMessage("validation failed").
			WithKind(richerror.KindInvalid).
			WithMetaMap(toAnyMap(fieldErrors))
	}

	return nil
}

func (v Validator) ValidateReorderPages(ctx context.Context, req param.ReorderPagesParam) error {
	const op = richerror.Op("validator.chapter.ValidateReorderPages")

	fieldErrors := make(map[string]string)

	if req.ChapterID == 0 {
		fieldErrors["chapter_id"] = "chapter ID is required"
	}

	if len(req.Pages) == 0 {
		fieldErrors["pages"] = "at least one page is required"
	}

	// Check for duplicate page numbers
	seen := make(map[int]bool)
	for i, p := range req.Pages {
		if p.PageID == 0 {
			fieldErrors[fmt.Sprintf("pages[%d].page_id", i)] = "page ID is required"
		}
		if p.PageNumber < 0 {
			fieldErrors[fmt.Sprintf("pages[%d].page_number", i)] = "page number must be non-negative"
		}
		if seen[p.PageNumber] {
			fieldErrors[fmt.Sprintf("pages[%d].page_number", i)] = "duplicate page number"
		}
		seen[p.PageNumber] = true
	}

	if len(fieldErrors) > 0 {
		return richerror.New(op).
			WithMessage("validation failed").
			WithKind(richerror.KindInvalid).
			WithMetaMap(toAnyMap(fieldErrors))
	}

	return nil
}

// ============================================
// Helpers
// ============================================

func validatePageFile(header *multipart.FileHeader, maxSizeMB int, allowedMimeTypes []string) error {
	if header.Size == 0 {
		return fmt.Errorf("file is empty")
	}

	maxSizeBytes := int64(maxSizeMB) * 1024 * 1024
	if header.Size > maxSizeBytes {
		return fmt.Errorf("file size exceeds %dMB limit", maxSizeMB)
	}

	mimeType := strings.ToLower(strings.TrimSpace(header.Header.Get("Content-Type")))
	allowed := false
	for _, m := range allowedMimeTypes {
		if strings.ToLower(m) == mimeType {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("invalid file type: %s", mimeType)
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".webp", ".gif"}
	extAllowed := false
	for _, e := range allowedExts {
		if e == ext {
			extAllowed = true
			break
		}
	}
	if !extAllowed {
		return fmt.Errorf("invalid file extension: %s", ext)
	}

	return nil
}

func toAnyMap(m map[string]string) map[string]any {
	result := make(map[string]any, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}