package domain

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type EchoProvider interface {
	GetEcho() *echo.Echo
}

type ConfigLoader interface {
	Load(name string, in interface{}) error
}

type SQLXProvider interface {
	GetSQLX(connectionString string) (*sqlx.DB, error)
}

type GatewayRepositoryProvider interface {
	GetRepository(db *sqlx.DB) GatewayRepository
}

type StorageV1ClientProvider interface {
	GetStorageV1Client() StorageV1Client
}

type GatewayUploaderServiceProvider interface {
	GetGatewayUploaderService(
		client StorageV1Client,
		repo GatewayRepository,
		concurrency int,
		logger *zap.Logger,
	) GatewayUploaderService
}

type GatewayDownloaderServiceProvider interface {
	GetGatewayDownloaderService(client StorageV1Client, logger *zap.Logger) GatewayDownloaderService
}

type GatewayDeleterServiceProvider interface {
	GetGatewayDeleterService(
		client StorageV1Client,
		repo GatewayRepository,
		concurrency int,
		logger *zap.Logger,
	) GatewayDeleterService
}

type GatewayServiceProvider interface {
	GetGatewayService(
		repo GatewayRepository,
		uploader GatewayUploaderService,
		downloader GatewayDownloaderService,
		deleter GatewayDeleterService,
		numChunks int64,
		minFileSizeToSplit int64,
		logger *zap.Logger,
	) GatewayService
}

type GatewayHandlerProvider interface {
	GetGatewayHandler(service GatewayService, logger *zap.Logger) EchoHandler
}

type StorageServiceProvider interface {
	GetStorageService(storagePath string) StorageService
}

type StorageHandlerProvider interface {
	GetStorageHandler(service StorageService, logger *zap.Logger) EchoHandler
}
