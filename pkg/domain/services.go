package domain

import (
	"context"
	"errors"
	"io"
)

var ErrNoStorageSpace = errors.New("no storage space")

type GatewayUploaderService interface {
	Upload(ctx context.Context, chunks Chunks, file File, reader io.Reader) error
}

type GatewayDownloaderService interface {
	Download(ctx context.Context, chunks Chunks, writer io.Writer) error
}

type GatewayDeleterService interface {
	Delete(ctx context.Context, chunks Chunks, fileID string) error
}

type GatewayService interface {
	UploadFile(ctx context.Context, file File, reader io.Reader) error
	DownloadFile(ctx context.Context, fileID string, writer io.Writer) error
	DeleteFile(ctx context.Context, fileID string) (File, error)

	AddServer(ctx context.Context, server Server) error
}

type StorageService interface {
	StoreChunk(ctx context.Context, fileName string, reader io.Reader) error
	GetChunk(ctx context.Context, fileName string, writer io.Writer) error
	DeleteChunk(ctx context.Context, name string) error
}
