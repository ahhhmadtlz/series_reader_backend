package imageprocessinghandler

import (
	"net/http"
	"strconv"

	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

// getCoverVariants handles GET /series/:seriesID/cover-variants
func (h Handler) getCoverVariants(c echo.Context) error {
	seriesID, err := parseID(c, "seriesID")
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid series ID",
		})
	}

	response, err := h.ipService.GetCoverVariants(c.Request().Context(), seriesID)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusOK, response)
}

// getBannerVariants handles GET /series/:seriesID/banner-variants
func (h Handler) getBannerVariants(c echo.Context) error {
	seriesID, err := parseID(c, "seriesID")
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid series ID",
		})
	}

	response, err := h.ipService.GetBannerVariants(c.Request().Context(), seriesID)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusOK, response)
}

// getThumbnailVariants handles GET /chapters/:chapterID/thumbnail-variants
func (h Handler) getThumbnailVariants(c echo.Context) error {
	chapterID, err := parseID(c, "chapterID")
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid chapter ID",
		})
	}

	response, err := h.ipService.GetThumbnailVariants(c.Request().Context(), chapterID)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusOK, response)
}

// getPageVariants handles GET /chapters/:chapterID/pages/:pageID/variants
func (h Handler) getPageVariants(c echo.Context) error {
	pageID, err := parseID(c, "pageID")
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid page ID",
		})
	}

	response, err := h.ipService.GetVariants(c.Request().Context(), pageID)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusOK, response)
}

func parseID(c echo.Context, paramName string) (uint, error) {
	val, err := strconv.ParseUint(c.Param(paramName), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}