package providers

import (
	"github.com/antelman107/filestorage/internal/services"
	"github.com/antelman107/filestorage/pkg/domain"
)

type defaultStorageServiceProvider struct {
}

func NewDefaultStorageServiceProvider() domain.StorageServiceProvider {
	return &defaultStorageServiceProvider{}
}

func (p *defaultStorageServiceProvider) GetStorageService(storagePath string) domain.StorageService {
	return services.NewStorageService(
		storagePath,
	)
}
