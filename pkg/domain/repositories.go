package domain

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("not found")

type GatewayRepository interface {
	LockServers(ctx context.Context) error
	StoreServer(ctx context.Context, server Server) error
	DeleteFile(ctx context.Context, id string) (File, error)
	GetServersUsages(ctx context.Context) (ServersSummary, error)

	StoreChunks(ctx context.Context, chunks Chunks) error
	UpdateUploadedChunk(ctx context.Context, id, hash string) error
	DeleteChunk(ctx context.Context, id string) error
	GetChunksByFileID(ctx context.Context, fileID string) (Chunks, error)

	LockFile(ctx context.Context, id string) error
	StoreFile(ctx context.Context, file File) error
	UpdateFileIsUploaded(ctx context.Context, id string) error
	GetUploadedFile(ctx context.Context, id string) (File, error)

	TransactionalRepository
}

type TransactionalRepository interface {
	InitTransaction() (*sqlx.Tx, error)
}
