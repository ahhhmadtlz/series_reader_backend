package storage

import (
	"context"
	"io"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/upload/entity"
)

type Storage interface {
	Save(ctx context.Context, req SaveRequest)(SaveResult,error)
	Delete(ctx context.Context, storedPath string)error
	URL(storedPath string)string
}


type SaveRequest struct {
	File io.Reader
	Filename string
	OwnerID uint
	Kind entity.ImageKind
	MimeType string
	Size int64
}


type SaveResult struct {
	StoredPath string
	URL string
}