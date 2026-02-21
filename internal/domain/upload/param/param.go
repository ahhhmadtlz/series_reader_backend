package param

import "mime/multipart"

// ============================================
// Upload Avatar
// ============================================

type UploadAvatarRequest struct {
	UserID uint
	File multipart.File
	Header *multipart.FileHeader
}

type UploadAvatarResponse struct {
	AvatarURL string `json:"avatar_url"`
}

// ============================================
// Upload Series Cover
// ============================================

type UploadCoverRequest struct {
	SeriesID uint // from URL path param
	UserID   uint // from JWT context (for permission check)
	File     multipart.File
	Header   *multipart.FileHeader
}

type UploadCoverResponse struct {
	CoverImageURL string `json:"cover_image_url"`
}

// ============================================
// Upload Chapter Page (TBD - scaffolded)
// ============================================

type UploadChapterPageRequest struct {
	ChapterID  uint
	PageNumber int
	UserID     uint // from JWT context
	File       multipart.File
	Header     *multipart.FileHeader
}

type UploadChapterPageResponse struct {
	PageImageURL string `json:"page_image_url"`
	PageNumber   int    `json:"page_number"`
}