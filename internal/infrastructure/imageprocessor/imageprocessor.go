package imageprocessor

import "github.com/ahhhmadtlz/series_reader_backend/internal/infrastructure/storage"

type ImageProcessor struct {
	storage  storage.Storage
	basePath string
}

func New(storage storage.Storage, basePath string) *ImageProcessor {
	return &ImageProcessor{
		storage:  storage,
		basePath: basePath,
	}
}