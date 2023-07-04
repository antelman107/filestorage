package providers

import (
	"github.com/antelman107/filestorage/internal/clients"
	"github.com/antelman107/filestorage/pkg/domain"
)

type defaultStorageClientProvider struct {
}

func NewDefaultStorageClientProvider() domain.StorageV1ClientProvider {
	return &defaultStorageClientProvider{}
}

func (p *defaultStorageClientProvider) GetStorageV1Client() domain.StorageV1Client {
	return clients.NewStorageClient()
}
