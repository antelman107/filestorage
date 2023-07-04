package providers

import (
	"go.uber.org/zap"

	"github.com/antelman107/filestorage/internal/handlers"
	"github.com/antelman107/filestorage/pkg/domain"
)

type defaultStorageHandlerProvider struct {
}

func NewDefaultStorageHandlerProvider() domain.StorageHandlerProvider {
	return &defaultStorageHandlerProvider{}
}

func (p *defaultStorageHandlerProvider) GetStorageHandler(
	service domain.StorageService,
	logger *zap.Logger,
) domain.EchoHandler {
	return handlers.NewStorageHandler(service, logger)
}
