package param

import "time"

type CreateBookmarkRequest struct {
	SeriesID uint `json:"series_id"`
}

type CreateBookmarkResponse struct {
	Bookmark BookmarkInfo `json:"bookmark"`
}

type BookmarkInfo struct {
	ID        uint       `json:"id"`
	SeriesID  uint       `json:"series_id"`
	Series    SeriesInfo `json:"series"`
	CreatedAt time.Time  `json:"created_at"`
}

type SeriesInfo struct {
	ID            uint     `json:"id"`
	Title         string   `json:"title"`
	FullSlug      string   `json:"full_slug"`
	CoverImageURL string   `json:"cover_image_url"`
	Type          string   `json:"type"`
	Status        string   `json:"status"`
	Genres        []string `json:"genres"`
}

type GetBookmarksResponse struct {
	Bookmarks []BookmarkInfo `json:"bookmarks"`
	Total     int            `json:"total"`
}

type DeleteBookmarkResponse struct {
	Message string `json:"message"`
}