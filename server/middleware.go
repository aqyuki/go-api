package server

import (
	"log/slog"

	"github.com/labstack/echo/v4"
)

// NewCustomLogger returns a customized middleware for logging
func NewCustomLogger(l *slog.Logger) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			l.Info(
				"Request",
				slog.String("method", c.Request().Method),
				slog.String("path", c.Request().URL.Path),
				slog.String("remote", c.Request().RemoteAddr),
				slog.String("user_agent", c.Request().UserAgent()),
			)
			return next(c)
		}
	}
}
