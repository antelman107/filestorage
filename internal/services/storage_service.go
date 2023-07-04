package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/antelman107/filestorage/pkg/domain"
)

const (
	defaultPermissions = 0755
	numSubdirectories  = 3
)

type storageFilesService struct {
	storagePath string
}

func NewStorageService(storagePath string) domain.StorageService {
	return &storageFilesService{
		storagePath: storagePath,
	}
}

func (s *storageFilesService) StoreChunk(_ context.Context, name string, reader io.Reader) error {
	fullPath := s.storagePath + string(os.PathSeparator) + getPathWithDirectories(name)

	directories := filepath.Dir(fullPath)
	if err := os.MkdirAll(directories, defaultPermissions); err != nil {
		return fmt.Errorf("failed to create directtory: %w, file %s, directories %s", err, name, directories)
	}

	destination, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, defaultPermissions)
	if err != nil {
		return fmt.Errorf("failed to open destination file: %w", err)
	}
	defer destination.Close()

	// Copy
	if _, err = io.Copy(destination, reader); err != nil {
		return fmt.Errorf("failed to copy chunk to destination file: %w", err)
	}

	return nil
}

func (s *storageFilesService) GetChunk(_ context.Context, name string, writer io.Writer) error {
	fullPath := s.storagePath + string(os.PathSeparator) + getPathWithDirectories(name)

	file, err := os.OpenFile(fullPath, os.O_RDONLY, defaultPermissions)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Copy
	if _, err = io.Copy(writer, file); err != nil {
		return fmt.Errorf("failed to copy to file: %w", err)
	}

	return nil
}

func (s *storageFilesService) DeleteChunk(_ context.Context, name string) error {
	fullPath := s.storagePath + string(os.PathSeparator) + getPathWithDirectories(name)

	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Delete subdirectories
	dir := fullPath
	for i := 0; i < numSubdirectories; i++ {
		dir = filepath.Dir(dir)
		if isDirectoryEmpty(dir) {
			if err := os.Remove(dir); err != nil {
				return fmt.Errorf("failed to delete directory %s: %w", dir, err)
			}
		}
	}

	return nil
}
