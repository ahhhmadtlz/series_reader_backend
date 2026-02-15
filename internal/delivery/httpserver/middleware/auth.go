package middleware

import (
	cfg "github.com/ahhhmadtlz/series_reader_backend/internal/config"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/auth"
	mw "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func Auth(service auth.Service, config auth.Config) echo.MiddlewareFunc {
	return mw.WithConfig(mw.Config{
		ContextKey:    cfg.AuthMiddlewareContextKey,
		SigningKey:    []byte(config.SignKey),
		SigningMethod: "HS256",
		ParseTokenFunc: func(c echo.Context, authToken string) (any, error) {
			claims, err := service.ParseBearerToken(authToken)
			if err != nil {
				return nil, err
			}
			return claims, nil
		},
	})
}