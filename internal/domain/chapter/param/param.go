package param

type CreateChapterRequest struct {
	SeriesID      uint    `json:"series_id"`
	ChapterNumber float64 `json:"chapter_number"`
	Title         *string `json:"title"`
}

type AddChapterPagesRequest struct {
	ChapterID uint                    `json:"chapter_id"`
	Pages     []CreateChapterPageItem `json:"pages"`
}

type CreateChapterPageItem struct {
	PageNumber int    `json:"page_number"`
	ImageURL   string `json:"image_url"`
}

type ChapterPageResponse struct {
	ID         uint   `json:"id"`
	PageNumber int    `json:"page_number"`
	ImageURL   string `json:"image_url"`
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
