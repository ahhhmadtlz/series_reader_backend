package middleware

import (
	cfg "github.com/ahhhmadtlz/series_reader_backend/internal/config"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/auth"
	"github.com/labstack/echo/v4"
)

func UserContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get(cfg.AuthMiddlewareContextKey).(*auth.Claims)
			if !ok {
				return echo.NewHTTPError(401, "unauthorized")
			}
			
			c.Set("user_id", claims.UserID)
			c.Set("user_role",claims.Role)
			c.Set("subscription_tier",claims.SubscriptionTier)
			
			return next(c)
		}
	}
}