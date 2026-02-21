package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/storage"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"

	"github.com/segmentio/ksuid"
)

type LocalStorage struct {
	basePath string
	baseURL  string
}

func New(basePath, baseURL string) *LocalStorage {
	return &LocalStorage{
		basePath: basePath,
		baseURL:  strings.TrimRight(baseURL, "/"),
	}
}

func (l *LocalStorage) Save(ctx context.Context, req storage.SaveRequest)(storage.SaveResult,error){
	const op=richerror.Op("storage.local.Save")

	
	uuid := ksuid.New().String()
	ext, err := extensionFromMime(req.MimeType)
	if err != nil {
		return storage.SaveResult{}, richerror.New(op).
			WithErr(err).
			WithMessage("unsupported mime type").
			WithKind(richerror.KindInvalid)
	}

  filename :=uuid +ext

	var relativePath string
	switch req.Kind {
		case entity.ImageKindAvatar :
				relativePath =fmt.Sprintf("avatars/%d/%s",req.OwnerID,filename)
		case entity.ImageKindCover:
			relativePath=fmt.Sprintf("covers/%d/%s",req.OwnerID,filename)
		case entity.ImageKindChapterPage:
			relativePath = fmt.Sprintf("chapters/%d/%s", req.OwnerID, filename)
		default:
			return storage.SaveResult{}, richerror.New(op).
				WithMessage("invalid image kind").
				WithKind(richerror.KindInvalid)
	}

	fullPath :=filepath.Join(l.basePath,relativePath)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return storage.SaveResult{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to create directory").
			WithKind(richerror.KindUnexpected)
	}

	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return storage.SaveResult{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to create file").
			WithKind(richerror.KindUnexpected)
	}
	defer file.Close()

	if _, err := io.Copy(file, req.File); err != nil {
		_ = os.Remove(fullPath) // Cleanup partial file on failure
		return storage.SaveResult{}, richerror.New(op).
			WithErr(err).
			WithMessage("failed to write file").
			WithKind(richerror.KindUnexpected)
	}


	return storage.SaveResult{
		StoredPath: relativePath,
		URL:        l.URL(relativePath),
	}, nil

}



func (l *LocalStorage) Delete(ctx context.Context, storedPath string) error {
	const op = richerror.Op("storage.local.Delete")

	fullPath := filepath.Join(l.basePath, storedPath)

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return richerror.New(op).
				WithMessage("file not found").
				WithKind(richerror.KindNotFound)
		}

		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to delete file").
			WithKind(richerror.KindUnexpected)
	}

	return nil
}


func (l *LocalStorage) URL(storedPath string) string {
	return l.baseURL + "/" + storedPath
}



func extensionFromMime(mime string) (string, error) {
	switch mime {
		case "image/jpeg", "image/jpg":
			return ".jpg", nil
		case "image/png":
			return ".png", nil
		case "image/webp":
			return ".webp", nil
		case "image/gif":
			return ".gif", nil
		default:
			return "", fmt.Errorf("unsupported mime type: %s", mime)
	}
}