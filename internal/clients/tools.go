package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/antelman107/filestorage/pkg/domain"
)

func composeURL(serverURL, urlPart string) string {
	return serverURL + urlPart
}

func getHTTPResponse(client http.Client, req *http.Request) (*http.Response, error) {
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do http request: %w", err)
	}

	if !isHTTPCodeOK(response.StatusCode) {
		_ = response.Body.Close()
		return nil, errors.New(fmt.Sprintf("http code is not OK: %d", response.StatusCode))
	}

	return response, nil
}

func getJSONHTTPResponse(client http.Client, req *http.Request, model interface{}) error {
	response, err := getHTTPResponse(client, req)
	if err != nil {
		return fmt.Errorf("failed to get http response: %w", err)
	}
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(model); err != nil {
		return fmt.Errorf("failed to decode json: %w", err)
	}

	return nil
}

func getHealth(ctx context.Context, serverURL string, httpClient http.Client) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/health", serverURL), nil)
	if err != nil {
		return fmt.Errorf("failed to build health request: %w", err)
	}

	_, err = getHTTPResponse(httpClient, req)

	return err
}

func isHTTPCodeOK(code int) bool {
	return code >= http.StatusOK && code <= 399
}

func getChunkUploadV1Request(ctx context.Context, chunk domain.ChunkWithData) (*http.Request, error) {
	body, writer, err := getChunkRequestBody(chunk)
	if err != nil {
		return nil, fmt.Errorf("failed to get chunk request body")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, composeURL(chunk.ServerURL, "/v1/chunks"), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create new post chunk request: %w", err)
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	return req, nil
}

func getPostFileUploadV1Request(ctx context.Context, name, serverURL string, readCloser io.ReadCloser) (*http.Request, error) {
	body, writer, err := getFileUploadRequestBody(name, readCloser)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, composeURL(serverURL, "/v1/files"), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create new post chunk request: %w", err)
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	return req, nil
}

func getFileUploadRequestBody(localFileName string, readCloser io.ReadCloser) (io.Reader, *multipart.Writer, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", localFileName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, readCloser); err != nil {
		return nil, nil, fmt.Errorf("failed to copy file files_for_upload to multipart writer: %w", err)
	}
	_ = readCloser.Close()

	if err := writer.Close(); err != nil {
		return nil, nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	return body, writer, nil
}

func getChunkRequestBody(chunk domain.ChunkWithData) (io.Reader, *multipart.Writer, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := writer.WriteField("fileName", chunk.ID.String()); err != nil {
		return nil, nil, fmt.Errorf("failed to write fileName field: %w", err)
	}
	if err := writer.WriteField("fileSize", fmt.Sprintf("%d", chunk.Size)); err != nil {
		return nil, nil, fmt.Errorf("failed to write fileSize field: %w", err)
	}

	part, err := writer.CreateFormFile("chunk", "chunk")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, bytes.NewReader(chunk.Data)); err != nil {
		return nil, nil, fmt.Errorf("failed to copy chunk files_for_upload to multipart writer: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	return body, writer, nil
}
