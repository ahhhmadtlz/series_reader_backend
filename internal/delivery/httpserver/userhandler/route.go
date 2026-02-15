package userhandler

import (
	"github.com/ahhhmadtlz/series_reader_backend/internal/delivery/httpserver/middleware"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/auth"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo,authService auth.Service,authConfig auth.Config) {
	publicGroup:= e.Group("/users")

	publicGroup.POST("/register", h.register)
	publicGroup.POST("/login", h.login)
	publicGroup.POST("/refresh", h.refresh)

	protectedGroup:=e.Group("/users")
	protectedGroup.Use(middleware.Auth(authService,authConfig))
	protectedGroup.Use(middleware.UserContext())

	protectedGroup.GET("/profile", h.getProfile)
  protectedGroup.PUT("/profile", h.updateProfile)
	protectedGroup.PUT("/password",h.changePassword)
}