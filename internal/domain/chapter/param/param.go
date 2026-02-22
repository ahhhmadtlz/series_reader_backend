package param

import "mime/multipart"

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


type ChapterPageResponse struct {
	ID         uint   `json:"id"`
	PageNumber int    `json:"page_number"`
	ImageURL   string `json:"image_url"`
}


type UploadPageParam struct {
	ChapterID uint
	PageNumber int
	File     multipart.File
	Header  *multipart.FileHeader
}

type BulkUploadParam struct {
	ChapterID uint
	Files []*multipart.FileHeader
}
type ReorderPageItem struct {
	PageID     uint `json:"page_id"`
	PageNumber int  `json:"page_number"`
}

type ReorderPagesParam struct {
	ChapterID uint
	Pages     []ReorderPageItem `json:"pages"`
}