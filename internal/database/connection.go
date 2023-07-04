package database

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
)

const driver = "pgx"

//go:embed gateway_migrations/*.sql
var gatewayMigrations embed.FS

// GetGatewayMigrations returns all embedded gatewayMigrations as embed.FS.
func GetGatewayMigrations() embed.FS { return gatewayMigrations }

func Get(connectionString string) (*sqlx.DB, error) {
	db, err := sql.Open(driver, connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to sqxl.Open: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to db.Ping: %w", err)
	}

	if err := applyMigrations(db, GetGatewayMigrations()); err != nil {
		return nil, fmt.Errorf("failed to apply gateway_migrations: %w", err)
	}

	return sqlx.NewDb(db, driver), nil
}

func applyMigrations(connection *sql.DB, fs embed.FS) error {
	mgs := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: fs,
		Root:       "gateway_migrations",
	}
	_, err := migrate.Exec(
		connection,
		"postgres",
		mgs,
		migrate.Up,
	)

	return err
}
