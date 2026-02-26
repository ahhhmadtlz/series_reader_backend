package bookmarkhandler

import (
	"context"

	bookmarkparam "github.com/ahhhmadtlz/series_reader_backend/internal/domain/bookmark/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/bookmark/validator"
)

type BookmarkService interface {
    CreateBookmark(ctx context.Context, userID uint, req bookmarkparam.CreateBookmarkRequest) (bookmarkparam.CreateBookmarkResponse, error)
    GetBookmarks(ctx context.Context, userID uint) (bookmarkparam.GetBookmarksResponse, error)
    DeleteBookmark(ctx context.Context, userID uint, seriesID uint) (bookmarkparam.DeleteBookmarkResponse, error)
}

type Handler struct {
    service   BookmarkService
    validator validator.Validator
}

func New(service BookmarkService, validator validator.Validator) Handler {
    return Handler{
        service:   service,
        validator: validator,
    }
}