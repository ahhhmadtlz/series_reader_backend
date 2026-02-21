package entity

import "time"

type UploadedImage struct {
	ID         uint
	OwnerID    uint // UserID for avatar, SeriesID for cover, ChapterID for page
	Kind       ImageKind
	Filename   string // original filename (sanitized)
	StoredPath string // relative path: "avatars/123/abc123.webp" or "covers/456/def456.webp"
	URL        string // public URL: "/uploads/avatars/123/abc123.webp"
	MimeType   string // "image/jpeg", "image/png", "image/webp"
	SizeBytes  int64
	CreatedAt  time.Time
}