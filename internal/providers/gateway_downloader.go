package providers

import (
	"go.uber.org/zap"

	"github.com/antelman107/filestorage/internal/services"
	"github.com/antelman107/filestorage/pkg/domain"
)

type defaultGatewayDownloaderProvider struct {
}

func NewDefaultGatewayDownloaderProvider() domain.GatewayDownloaderServiceProvider {
	return &defaultGatewayDownloaderProvider{}
}

func (p *defaultGatewayDownloaderProvider) GetGatewayDownloaderService(
	client domain.StorageV1Client,
	logger *zap.Logger,
) domain.GatewayDownloaderService {
	return services.NewGatewayDownloader(client, logger)
}
