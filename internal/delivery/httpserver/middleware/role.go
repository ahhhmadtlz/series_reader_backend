package middleware

import (
	"net/http"

	sharedentity "github.com/ahhhmadtlz/series_reader_backend/internal/domain/shared/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func RequireRole(allowedRoles ...sharedentity.Role) echo.MiddlewareFunc {
	return  func(next echo.HandlerFunc) echo.HandlerFunc{
		return  func (c echo.Context)error {
			userRole, ok :=c.Get("user_role").(sharedentity.Role)
			if !ok {
				return c.JSON(http.StatusUnauthorized,httpmsgerrorhandler.ErrorResponse{
					Message: "unauthorized",
				})
			}

			for _,role :=range allowedRoles {
				if userRole ==role {
					return next(c)
				}
			}
			return c.JSON(http.StatusForbidden, httpmsgerrorhandler.ErrorResponse{
				Message: "forbidden - insufficient role",
			})
		}
	}
}


func RequireAdmin() echo.MiddlewareFunc {
	return RequireRole(sharedentity.AdminRole)
}

func RequireManagerOrAdmin() echo.MiddlewareFunc {
	return RequireRole(sharedentity.AdminRole, sharedentity.ManagerRole)
}


func RequirePremium() echo.MiddlewareFunc {
	return  func(next echo.HandlerFunc)echo.HandlerFunc {
		return func(c echo.Context) error{
			tier, ok := c.Get("subscription_tier").(sharedentity.SubscriptionTier)
			if !ok {
				return  c.JSON(http.StatusUnauthorized,httpmsgerrorhandler.ErrorResponse{
					Message: "unauthorized",
				})
			}

			if tier !=sharedentity.PremiumTier{
				return c.JSON(http.StatusForbidden, httpmsgerrorhandler.ErrorResponse{
					Message: "forbidden - premium subscription required",
				})
			}
			return next(c)
		}
	}
}