package providers

import (
	"go.uber.org/zap"

	"github.com/antelman107/filestorage/internal/services"
	"github.com/antelman107/filestorage/pkg/domain"
)

type defaultGatewayServiceProvider struct {
}

func NewDefaultGatewayServiceProvider() domain.GatewayServiceProvider {
	return &defaultGatewayServiceProvider{}
}

func (p *defaultGatewayServiceProvider) GetGatewayService(
	repo domain.GatewayRepository,
	uploader domain.GatewayUploaderService,
	downloader domain.GatewayDownloaderService,
	deleter domain.GatewayDeleterService,
	numChunks int64,
	minFileSizeToSplit int64,
	logger *zap.Logger,
) domain.GatewayService {
	return services.NewGatewayService(
		repo,
		uploader, downloader, deleter,
		numChunks, minFileSizeToSplit,
		logger,
	)
}
