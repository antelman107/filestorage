package providers

import (
	"go.uber.org/zap"

	"github.com/antelman107/filestorage/internal/services"
	"github.com/antelman107/filestorage/pkg/domain"
)

type defaultGatewayUploaderProvider struct {
}

func NewDefaultGatewayUploaderProvider() domain.GatewayUploaderServiceProvider {
	return &defaultGatewayUploaderProvider{}
}

func (p *defaultGatewayUploaderProvider) GetGatewayUploaderService(
	client domain.StorageV1Client,
	repo domain.GatewayRepository,
	concurrency int,
	logger *zap.Logger,
) domain.GatewayUploaderService {
	return services.NewGatewayUploader(client, repo, concurrency, logger)
}
