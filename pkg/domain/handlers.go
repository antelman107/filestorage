package domain

import "github.com/labstack/echo/v4"

type EchoHandler interface {
	AssignHandlers(e *echo.Echo)
}
