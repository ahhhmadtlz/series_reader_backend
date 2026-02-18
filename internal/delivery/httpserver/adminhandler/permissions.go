package adminhandler

import (
	"net/http"
	"strconv"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/user/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) getUserPermissions(c echo.Context) error {
	targetUserIDStr := c.Param("user_id")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid user_id",
		})
	}

	response, err := h.userService.GetUserWithPermissions(c.Request().Context(), uint(targetUserID))
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusOK, response)
}

func (h Handler) grantPermission(c echo.Context) error {
	adminID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpmsgerrorhandler.ErrorResponse{
			Message: "unauthorized",
		})
	}

	targetUserIDStr := c.Param("user_id")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid user_id",
		})
	}

	var req param.GrantPermissionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid request body",
		})
	}

	response, err := h.userService.GrantPermission(c.Request().Context(), uint(targetUserID), req, adminID)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusOK, response)
}

func (h Handler) revokePermission(c echo.Context) error {
	targetUserIDStr := c.Param("user_id")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid user_id",
		})
	}

	var req param.RevokePermissionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpmsgerrorhandler.ErrorResponse{
			Message: "invalid request body",
		})
	}

	response, err := h.userService.RevokePermission(c.Request().Context(), uint(targetUserID), req)
	if err != nil {
		return httpmsgerrorhandler.Error(c, err)
	}

	return c.JSON(http.StatusOK, response)
}