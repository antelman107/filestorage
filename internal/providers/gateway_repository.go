package providers

import (
	"github.com/jmoiron/sqlx"

	"github.com/antelman107/filestorage/internal/database/repositories"
	"github.com/antelman107/filestorage/pkg/domain"
)

type defaultGatewayRepoProvider struct {
}

func NewDefaultGatewayRepositoryProvider() domain.GatewayRepositoryProvider {
	return &defaultGatewayRepoProvider{}
}

func (p *defaultGatewayRepoProvider) GetRepository(db *sqlx.DB) domain.GatewayRepository {
	return repositories.NewGatewayRepo(db)
}
