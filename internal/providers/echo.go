package providers

import (
	"github.com/labstack/echo/v4"

	"github.com/antelman107/filestorage/pkg/domain"
)

type defaultEchoProvider struct {
}

func NewDefaultEchoProvider() domain.EchoProvider {
	return &defaultEchoProvider{}
}

func (p *defaultEchoProvider) GetEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	return e
}
