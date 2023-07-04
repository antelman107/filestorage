package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/antelman107/filestorage/pkg/domain"
)

type storageHandler struct {
	service domain.StorageService
	logger  *zap.Logger
}

func NewStorageHandler(service domain.StorageService, logger *zap.Logger) domain.EchoHandler {
	return &storageHandler{service: service, logger: logger}
}

func (h *storageHandler) AssignHandlers(e *echo.Echo) {
	e.POST("/v1/chunks", h.post)
	e.GET("/v1/chunks/:chunk_id", h.get)
	e.DELETE("/v1/chunks/:chunk_id", h.delete)
	e.GET("/health", healthHandler)
}

type postRequest struct {
	fileName   string `validate:"required"`
	fileSize   int64  `validate:"required"`
	fileReader multipart.File
}

func (h *storageHandler) post(c echo.Context) error {
	req, err := decodePostRequest(c)
	if err != nil {
		return logHTTPErr(err, http.StatusBadRequest, "failed to decode postFiles request", h.logger)
	}
	defer req.fileReader.Close()

	if err := h.service.StoreChunk(c.Request().Context(), req.fileName, req.fileReader); err != nil {
		return logHTTPErr(err, http.StatusBadRequest, "failed to store chunk", h.logger)
	}

	return c.NoContent(http.StatusOK)
}

func (h *storageHandler) get(c echo.Context) error {
	chunkID := c.Param("chunk_id")
	if _, err := uuid.Parse(chunkID); err != nil {
		return logHTTPErr(err, http.StatusBadRequest, "failed to parse chunk_id UUID", h.logger)
	}

	if err := h.service.GetChunk(c.Request().Context(), chunkID, c.Response()); err != nil {
		return logHTTPErr(err, http.StatusBadRequest, "failed to getContent chunk", h.logger)
	}

	return nil
}

func (h *storageHandler) delete(c echo.Context) error {
	chunkID := c.Param("chunk_id")
	if _, err := uuid.Parse(chunkID); err != nil {
		return logHTTPErr(err, http.StatusBadRequest, "failed to parse chunk_id UUID", h.logger)
	}

	if err := h.service.DeleteChunk(c.Request().Context(), chunkID); err != nil {
		return logHTTPErr(err, http.StatusBadRequest, "failed to delete chunk", h.logger)
	}

	return nil
}

func decodePostRequest(c echo.Context) (*postRequest, error) {
	res := postRequest{}

	values, err := c.FormParams()
	if err != nil {
		return nil, fmt.Errorf("failed to getContent form params: %w", err)
	}

	res.fileName = values.Get("fileName")

	fileSize := values.Get("fileSize")
	fileSizeInt, err := strconv.ParseInt(fileSize, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode fileSize: %w", err)
	}

	res.fileSize = fileSizeInt

	multiPartFile, err := c.FormFile("chunk")
	if err != nil {
		return nil, fmt.Errorf("failed to getContent form file: %w", err)
	}

	source, err := multiPartFile.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	res.fileReader = source

	if err := validator.New().Struct(res); err != nil {
		return nil, fmt.Errorf("failed to validate postFiles request: %w", err)
	}

	return &res, nil
}
