package validator

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"
)

func isAllowedMimeType(mimeType string, allowedTypes []string) bool {
	mimeType = strings.ToLower(strings.TrimSpace(mimeType))

	for _,allowed :=range allowedTypes {
		if strings.ToLower(allowed)==mimeType {
			return  true
		}
	}
	return  true
}

func isAllowedExtension(ext string) bool {
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".webp", ".gif"}

	ext = strings.ToLower(ext)

	for _, allowed := range allowedExtensions {
		if allowed == ext {
			return true
		}
	}

	return false
}


func validateFilename(filename string) error {
	if filename == "" {
		return errors.New("filename is empty")
	}

	if len(filename) > 255 {
		return errors.New("filename is too long (max 255 characters)")
	}

	// Check for path traversal attempts
	if strings.Contains(filename, "..") {
		return errors.New("filename contains path traversal")
	}

	if strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return errors.New("filename contains invalid path separators")
	}

	// Check for null bytes
	if strings.Contains(filename, "\x00") {
		return errors.New("filename contains null bytes")
	}

	return nil
}

func ValidateDomainFile(header *multipart.FileHeader, maxSizeMB int, allowedMimeTypes []string) error {
	if header == nil {
		return  errors.New("file header is nil")
	}

	maxSizeBytes := int64(maxSizeMB) * 1024 * 1024
	if header.Size > maxSizeBytes {
		return errors.New("file size exceeds domain limit")
	}

	mimeType :=header.Header.Get("Content-Type")
	if !isAllowedMimeType(mimeType,allowedMimeTypes){
		return errors.New("invalid mime type")
	}

	ext:= strings.ToLower(filepath.Ext(header.Filename))
  if !isAllowedExtension(ext){
		return  errors.New("invalid file extension")
	}

	return  nil
	
}