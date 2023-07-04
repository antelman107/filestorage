package clients

import (
	"context"
	"net/http"

	"github.com/antelman107/filestorage/pkg/domain"
)

type healthClient struct {
	http.Client
}

func NewHealthClient() domain.HealthClient {
	return &healthClient{}
}

func (c *healthClient) GetHealth(ctx context.Context, serverURL string) error {
	return getHealth(ctx, serverURL, c.Client)
}
