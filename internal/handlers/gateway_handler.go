package handlers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/antelman107/filestorage/pkg/domain"
)

type gatewayHandler struct {
	service domain.GatewayService
	logger  *zap.Logger
}

func NewGatewayHandler(service domain.GatewayService, logger *zap.Logger) domain.EchoHandler {
	return &gatewayHandler{
		service: service,
		logger:  logger,
	}
}

func (h *gatewayHandler) AssignHandlers(e *echo.Echo) {
	e.POST("/v1/files", h.postFiles)
	e.GET("/v1/files/:file_id/content", h.getContent)
	e.DELETE("/v1/files/:file_id", h.deleteFile)

	e.POST("/v1/servers", h.postServer)

	e.GET("/health", healthHandler)
}

func (h *gatewayHandler) postFiles(c echo.Context) error {
	requestFile, err := c.FormFile("file")
	if err != nil {
		return logHTTPErr(err, http.StatusBadRequest, "failed to parse file parameter", h.logger)
	}

	source, err := requestFile.Open()
	if err != nil {
		return logHTTPErr(err, http.StatusInternalServerError, "failed to open file", h.logger)
	}
	defer source.Close()

	file := domain.File{
		ID:   uuid.New(),
		Name: requestFile.Filename,
		Size: requestFile.Size,
	}

	if err := h.service.UploadFile(c.Request().Context(), file, source); err != nil {
		if errors.Is(err, domain.ErrNoStorageSpace) {
			return logHTTPErr(err, http.StatusPreconditionFailed, "no storage space", h.logger)
		}
		return logHTTPErr(err, http.StatusInternalServerError, "failed to upload file", h.logger)
	}

	return c.JSON(http.StatusCreated, file)
}

func (h *gatewayHandler) getContent(c echo.Context) error {
	fileID := c.Param("file_id")

	if err := h.service.DownloadFile(c.Request().Context(), fileID, c.Response().Writer); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return logHTTPErr(err, http.StatusNotFound, "file is not found", h.logger)
		}

		return logHTTPErr(err, http.StatusInternalServerError, "failed to download file", h.logger)
	}

	return nil
}

type postServerRequest struct {
	URL string `json:"url"`
}

func (h *gatewayHandler) postServer(c echo.Context) error {
	var req postServerRequest
	err := c.Bind(&req)
	if err != nil {
		return logHTTPErr(err, http.StatusBadRequest, "failed to parse post server request", h.logger)
	}

	server := domain.Server{
		ID:  uuid.New(),
		URL: req.URL,
	}

	if err := h.service.AddServer(c.Request().Context(), server); err != nil {
		return logHTTPErr(err, http.StatusInternalServerError, "failed to add server", h.logger)
	}

	return c.JSON(http.StatusCreated, server)
}

func (h *gatewayHandler) deleteFile(c echo.Context) error {
	file, err := h.service.DeleteFile(c.Request().Context(), c.Param("file_id"))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return logHTTPErr(err, http.StatusNotFound, "file is not found", h.logger)
		}
		return logHTTPErr(err, http.StatusInternalServerError, "failed to delete file", h.logger)
	}

	return c.JSON(http.StatusCreated, file)
}
