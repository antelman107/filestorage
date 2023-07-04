package clients

import (
	"context"
	"fmt"
	"net/http"

	"github.com/antelman107/filestorage/pkg/domain"
)

type storageV1Client struct {
	http.Client
	healthClient
}

func NewStorageClient() domain.StorageV1Client {
	return &storageV1Client{}
}

func (c *storageV1Client) PostChunk(ctx context.Context, chunk domain.ChunkWithData) error {
	req, err := getChunkUploadV1Request(ctx, chunk)
	if err != nil {
		return fmt.Errorf("failed to build upload request: %w", err)
	}

	_, err = getHTTPResponse(c.Client, req)

	return err
}

func (c *storageV1Client) GetChunk(ctx context.Context, chunk domain.Chunk) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/v1/chunks/%s", chunk.ServerURL, chunk.ID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new get chunk request: %w", err)
	}

	return getHTTPResponse(c.Client, req)
}

func (c *storageV1Client) DeleteChunk(ctx context.Context, chunk domain.Chunk) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/v1/chunks/%s", chunk.ServerURL, chunk.ID), nil)
	if err != nil {
		return fmt.Errorf("failed to create new delete chunk request: %w", err)
	}

	response, err := getHTTPResponse(c.Client, req)
	defer response.Body.Close()

	return err
}
