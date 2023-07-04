//go:build integration
// +build integration

package integration_tests

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/antelman107/filestorage/internal/clients"
	"github.com/antelman107/filestorage/internal/config"
	"github.com/antelman107/filestorage/internal/providers"
)

type customSQLXProvider struct {
	db *sqlx.DB
}

func (c *customSQLXProvider) GetSQLX(connectionString string) (*sqlx.DB, error) {
	db, err := providers.NewDefaultSQLXProvider().GetSQLX(connectionString)
	c.db = db

	return db, err
}

type customStorageConfigLoader struct {
	config config.StorageConfig
}

func (c *customStorageConfigLoader) Load(_ string, in interface{}) error {
	cfg := in.(*config.StorageConfig)
	*cfg = c.config
	return nil
}

func isURLHealthy(ctx context.Context, url string) bool {
	client := clients.NewHealthClient()

	ticker := time.NewTicker(50 * time.Millisecond)
	timeoutCtx, cancelFunc := context.WithTimeout(ctx, time.Second*30)
	defer cancelFunc()

forloop:
	for {
		select {
		case <-timeoutCtx.Done():
			break forloop

		case <-ticker.C:
			err := client.GetHealth(timeoutCtx, url)
			if err == nil {
				return true
			}
		}

	}

	return false
}

func getDataFile(fileName string) (*os.File, error) {
	return os.Open("./files_for_upload/" + fileName)
}

func clearStorage() error {
	_ = os.RemoveAll("./storage/storage_data1/")
	_ = os.RemoveAll("./storage/storage_data2/")
	return nil
}

func clearDB(db *sqlx.DB) error {
	if _, err := db.Exec("DELETE FROM chunks;"); err != nil {
		return fmt.Errorf("failed to delete from chunks: %w", err)
	}
	if _, err := db.Exec("DELETE FROM files;"); err != nil {
		return fmt.Errorf("failed to delete from files: %w", err)
	}
	if _, err := db.Exec("DELETE FROM servers;"); err != nil {
		return fmt.Errorf("failed to delete from files: %w", err)
	}
	return nil
}

func getStorageFilesPaths() ([]string, error) {
	result := make([]string, 0)

	if err := filepath.Walk("./storage", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to iterate over files: %w", err)
		}
		if !info.IsDir() {
			result = append(result, path)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to iterate over files: %w", err)
	}

	return result, nil
}

func createDataFile(fileName string, size int, char byte) error {
	f, err := os.OpenFile("./files_for_upload/"+fileName, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", fileName, err)
	}
	defer f.Close()

	for i := 0; i < size; i++ {
		if _, err := f.Write([]byte{char}); err != nil {
			return fmt.Errorf("failed to write byte to file %s: %w", fileName, err)
		}
	}

	return nil
}
