package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func logHTTPErr(err error, code int, message string, logger *zap.Logger) *echo.HTTPError {
	logger.Error(message, zap.Error(err))
	return echo.NewHTTPError(code, message)
}

func healthHandler(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
