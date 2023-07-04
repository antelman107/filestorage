package domain

import (
	"context"
	"io"
	"net/http"
)

type HealthClient interface {
	GetHealth(ctx context.Context, serverURL string) error
}

type GatewayV1Client interface {
	PostFile(ctx context.Context, name string, file io.ReadCloser) (File, error)
	DeleteFile(ctx context.Context, id string) (File, error)
	GetFileContent(ctx context.Context, id string, writer io.Writer) error
	PostServer(ctx context.Context, url string) (Server, error)
	HealthClient
}

type StorageV1Client interface {
	PostChunk(ctx context.Context, chunk ChunkWithData) error
	GetChunk(ctx context.Context, chunk Chunk) (*http.Response, error)
	DeleteChunk(ctx context.Context, chunk Chunk) error
	HealthClient
}
