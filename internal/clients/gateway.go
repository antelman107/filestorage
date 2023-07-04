package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/antelman107/filestorage/pkg/domain"
)

type gatewayV1Client struct {
	http.Client
	healthClient

	serverURL string
}

func NewGatewayV1Client(serverURL string) domain.GatewayV1Client {
	return &gatewayV1Client{
		serverURL: serverURL,
	}
}

func (c *gatewayV1Client) PostFile(ctx context.Context, name string, readCloser io.ReadCloser) (domain.File, error) {
	req, err := getPostFileUploadV1Request(ctx, name, c.serverURL, readCloser)
	if err != nil {
		return domain.File{}, fmt.Errorf("failed to build upload request: %w", err)
	}

	file := domain.File{}
	if err := getJSONHTTPResponse(c.Client, req, &file); err != nil {
		return domain.File{}, fmt.Errorf("failed to get json http response: %w", err)
	}

	return file, nil
}

func (c *gatewayV1Client) GetFileContent(ctx context.Context, id string, writer io.Writer) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		composeURL(c.serverURL, fmt.Sprintf("/v1/files/%s/content", id)),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create get file request: %w", err)
	}

	response, err := getHTTPResponse(c.Client, req)
	if err != nil {
		return fmt.Errorf("failed to get http response: %w", err)
	}
	defer response.Body.Close()

	if _, err := io.Copy(writer, response.Body); err != nil {
		return fmt.Errorf("failed to copy response to writer: %w", err)
	}

	return nil
}

func (c *gatewayV1Client) DeleteFile(ctx context.Context, id string) (domain.File, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		composeURL(c.serverURL, fmt.Sprintf("/v1/files/%s", id)),
		nil,
	)
	if err != nil {
		return domain.File{}, fmt.Errorf("failed to create delete file request: %w", err)
	}

	res := domain.File{}
	if err := getJSONHTTPResponse(c.Client, req, &res); err != nil {
		return res, fmt.Errorf("failed to get http response: %w", err)
	}

	return res, nil
}

type postServerRequest struct {
	URL string `json:"URL"`
}

func (c *gatewayV1Client) PostServer(ctx context.Context, url string) (domain.Server, error) {
	js, err := json.Marshal(postServerRequest{URL: url})
	if err != nil {
		return domain.Server{}, fmt.Errorf("failed marshal url json: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		composeURL(c.serverURL, "/v1/servers"),
		bytes.NewReader(js),
	)
	if err != nil {
		return domain.Server{}, fmt.Errorf("failed to create get file request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	res := domain.Server{}
	if err := getJSONHTTPResponse(c.Client, req, &res); err != nil {
		return res, fmt.Errorf("failed to get http response: %w", err)
	}

	return res, nil
}
