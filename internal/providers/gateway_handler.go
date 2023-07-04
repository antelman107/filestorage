package providers

import (
	"go.uber.org/zap"

	"github.com/antelman107/filestorage/internal/handlers"
	"github.com/antelman107/filestorage/pkg/domain"
)

type defaultGatewayHandlerProvider struct {
}

func NewDefaultGatewayHandlerProvider() domain.GatewayHandlerProvider {
	return &defaultGatewayHandlerProvider{}
}

func (p *defaultGatewayHandlerProvider) GetGatewayHandler(
	service domain.GatewayService,
	logger *zap.Logger,
) domain.EchoHandler {
	return handlers.NewGatewayHandler(service, logger)
}
