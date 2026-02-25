package param

import (
	"mime/multipart"

	ipentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/entity"
	sharedentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/shared/entity"
)

type CreateChapterRequest struct {
	SeriesID      uint    `json:"series_id"`
	ChapterNumber float64 `json:"chapter_number"`
	Title         *string `json:"title"`
}

type ChapterResponse struct {
	ID            uint    `json:"id"`
	SeriesID      uint    `json:"series_id"`
	ChapterNumber float64 `json:"chapter_number"`
	Title         *string `json:"title"`
}

type ChapterWithPagesResponse struct {
	ID            uint                  `json:"id"`
	SeriesID      uint                  `json:"series_id"`
	ChapterNumber float64               `json:"chapter_number"`
	Title         *string               `json:"title"`
	Pages         []ChapterPageResponse `json:"pages"`
}

type PageVariantResponse struct {
	ID        uint                 `json:"id"`
	Kind      ipentity.VariantKind `json:"kind"`
	ImageURL  string               `json:"image_url"`
	CreatedAt string               `json:"created_at"`
}

type ChapterPageResponse struct {
	ID         uint                  `json:"id"`
	PageNumber int                   `json:"page_number"`
	ImageURL   string                `json:"image_url"`
	Variants   []PageVariantResponse `json:"variants"`
}

type UploadPageParam struct {
	ChapterID  uint
	PageNumber int
	File       multipart.File
	Header     *multipart.FileHeader
}

type BulkUploadParam struct {
	ChapterID  uint
	Files      []*multipart.FileHeader
	Force      bool
	CallerRole sharedentity.Role
}

type ReorderPageItem struct {
	PageID     uint `json:"page_id"`
	PageNumber int  `json:"page_number"`
}

type ReorderPagesParam struct {
	ChapterID uint
	Pages     []ReorderPageItem `json:"pages"`
}