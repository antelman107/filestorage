package services

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/antelman107/filestorage/internal/config"
	"github.com/antelman107/filestorage/internal/logger"
	"github.com/antelman107/filestorage/pkg/domain"
)

type storageApp struct {
	providers  StorageAppProviders
	echo       *echo.Echo
	listenPort string
	logger     *zap.Logger
}

type StorageAppProviders struct {
	EchoProvider    domain.EchoProvider
	ConfigLoader    domain.ConfigLoader
	ServiceProvider domain.StorageServiceProvider
	HandlerProvider domain.StorageHandlerProvider
}

func NewStorageApp(providers StorageAppProviders) domain.App {
	return &storageApp{providers: providers}
}

func (a *storageApp) Init() error {
	zapLogger, err := logger.Get()
	if err != nil {
		a.echo.Logger.Error("failed to init logger", err)
	}
	a.logger = zapLogger

	a.echo = a.providers.EchoProvider.GetEcho()

	cfg := config.DefaultStorageConfig
	if err := a.providers.ConfigLoader.Load("storage", &cfg); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	service := a.providers.ServiceProvider.GetStorageService(cfg.StoragePath)
	handler := a.providers.HandlerProvider.GetStorageHandler(service, zapLogger)
	handler.AssignHandlers(a.echo)

	a.listenPort = cfg.HTTP.ListenPort

	return nil
}

func (a *storageApp) Run(ctx context.Context) error {
	return runEchoWithGracefulShutdown(ctx, a.echo, a.listenPort, a.logger)
}
