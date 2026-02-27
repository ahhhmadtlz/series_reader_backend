package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (s Server) healthCheck(c echo.Context) error {
	// Ping DB with a short timeout so a slow DB doesn't hang the health check
	ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {
		return c.JSON(http.StatusServiceUnavailable, echo.Map{
			"status":  "unhealthy",
			"message": "database unreachable",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":  "healthy",
		"message": "everything is good ! keep going !",
	})
}