package middleware

import (
	"context"
	"net/http"

	sharedentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/shared/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

type PermissionChecker interface {
	HasPermission(ctx context.Context, userID uint, permission string) (bool, error)
}

func RequirePermission(checker PermissionChecker, permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, ok:=c.Get("user_id").(uint)
			if !ok{
				return c.JSON(http.StatusUnauthorized,httpmsgerrorhandler.ErrorResponse{
					Message: "unauthorized",
				})
			}
		  userRole, ok := c.Get("user_role").(sharedentity.Role)
			if !ok {
				return c.JSON(http.StatusUnauthorized, httpmsgerrorhandler.ErrorResponse{
					Message: "unauthorized",
				})
			}

			if userRole == sharedentity.AdminRole {
				return next(c)
			}
      hasPermission,err:=checker.HasPermission(c.Request().Context(),userID,permission)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, httpmsgerrorhandler.ErrorResponse{
					Message: "failed to check permissions",
				})
			}

			if !hasPermission {
					return c.JSON(http.StatusForbidden, httpmsgerrorhandler.ErrorResponse{
						Message: "forbidden - missing permission: " + permission,
					})
				}

				return next(c)
			}
		}
	}
