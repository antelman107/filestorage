package providers

import (
	"github.com/jmoiron/sqlx"

	"github.com/antelman107/filestorage/internal/database"
	"github.com/antelman107/filestorage/pkg/domain"
)

type defaultSQLXProvider struct {
}

func NewDefaultSQLXProvider() domain.SQLXProvider {
	return &defaultSQLXProvider{}
}

func (p *defaultSQLXProvider) GetSQLX(connectionString string) (*sqlx.DB, error) {
	return database.Get(connectionString)
}
