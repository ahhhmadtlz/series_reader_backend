package param

type GetListFilter struct {
	Status      string   `query:"status"`
	Type        string   `query:"type"`
	Genres      []string `query:"genres"`
	IsPublished *bool    `query:"is_published"`
	Search      string   `query:"search"` //Search in title , author, artist
}

type GetListSort struct {
	SortBy    string `query:"sort_by"`
	SortOrder string `query:"sort_order"`
}

type GetListPagination struct {
	Page     int `query:"page"`
	PageSize int `query:"page_size"`
}

type GetListRequest struct {
	Filter     GetListFilter
	Sort       GetListSort
	Pagination GetListPagination
}

type GetListResponse struct {
	Items      []SeriesResponse `json:"items"`
	Pagination PaginationMeta   `json:"pagination"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type CreateSeriesRequest struct {
	Title             string   `json:"title"`
	Slug              string   `json:"slug"`
	SlugID            string   `json:"slug_id"`
	FullSlug          string   `json:"full_slug"`
	Description       string   `json:"description"`
	Author            string   `json:"author"`
	Artist            string   `json:"artist"`
	Status            string   `json:"status"`
	Type              string   `json:"type"`
	Genres            []string `json:"genres"`
	AlternativeTitles []string `json:"alternative_titles"`
	CoverImageURL     string   `json:"cover_image_url"`
	PublicationYear   *int     `json:"publication_year"`
	IsPublished       bool     `json:"is_published"`
}

type UpdateSeriesRequest struct {
	Title             *string  `json:"title"`
	Description       *string  `json:"description"`
	Author            *string  `json:"author"`
	Artist            *string  `json:"artist"`
	Status            *string  `json:"status"`
	Type              *string  `json:"type"`
	Genres            []string `json:"genres"`
	AlternativeTitles []string `json:"alternative_titles"`
	CoverImageURL     *string  `json:"cover_image_url"`
	PublicationYear   *int     `json:"publication_year"`
	IsPublished       *bool    `json:"is_published"`
}

type SeriesResponse struct {
	ID                uint     `json:"id"`
	Title             string   `json:"title"`
	Slug              string   `json:"slug"`
	SlugID					  string   `json:"slug_id"`
	FullSlug          string   `json:"full_slug"`
	Description       string   `json:"description"`
	Author            string   `json:"author"`
	Artist            string   `json:"artist"`
	Status            string   `json:"status"`
	Type              string   `json:"type"`
	Genres            []string `json:"genres"`
	AlternativeTitles []string `json:"alternative_titles"`
	CoverImageURL     string   `json:"cover_image_url"`
	PublicationYear   *int     `json:"publication_year"`
	ViewCount         int      `json:"view_count"`
	Rating            float64  `json:"rating"`
	IsPublished       bool     `json:"is_published"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
}
