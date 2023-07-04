package providers

import (
	"go.uber.org/zap"

	"github.com/antelman107/filestorage/internal/services"
	"github.com/antelman107/filestorage/pkg/domain"
)

type defaultGatewayDeleterProvider struct {
}

func NewDefaultGatewayDeleterProvider() domain.GatewayDeleterServiceProvider {
	return &defaultGatewayDeleterProvider{}
}

func (p *defaultGatewayDeleterProvider) GetGatewayDeleterService(
	client domain.StorageV1Client,
	repo domain.GatewayRepository,
	concurrency int,
	logger *zap.Logger,
) domain.GatewayDeleterService {
	return services.NewGatewayDeleter(client, repo, concurrency, logger)
}
