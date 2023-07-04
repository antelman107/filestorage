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

// GatewayAppProviders contains all dependency providers for gateway app.
type GatewayAppProviders struct {
	EchoProvider            domain.EchoProvider
	ConfigLoader            domain.ConfigLoader
	SqlxProvider            domain.SQLXProvider
	RepoProvider            domain.GatewayRepositoryProvider
	StorageV1ClientProvider domain.StorageV1ClientProvider
	UploaderProvider        domain.GatewayUploaderServiceProvider
	DownloaderProvider      domain.GatewayDownloaderServiceProvider
	DeleterProvider         domain.GatewayDeleterServiceProvider
	ServiceProvider         domain.GatewayServiceProvider
	HandlerProvider         domain.GatewayHandlerProvider
}

type gatewayApp struct {
	providers  GatewayAppProviders
	echo       *echo.Echo
	listenPort string
	logger     *zap.Logger
}

func NewGatewayApp(providers GatewayAppProviders) domain.App {
	return &gatewayApp{
		providers: providers,
	}
}

// Init uses providers to resolve gateway app dependencies and store necessary fields.
func (a *gatewayApp) Init() error {
	zapLogger, err := logger.Get()
	if err != nil {
		a.echo.Logger.Error("failed to init logger", err)
	}
	a.logger = zapLogger

	a.echo = a.providers.EchoProvider.GetEcho()

	cfg := config.DefaultGatewayConfig
	if err := a.providers.ConfigLoader.Load("gateway", &cfg); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	db, err := a.providers.SqlxProvider.GetSQLX(cfg.DB.GetConnectionString())
	if err != nil {
		return fmt.Errorf("failed to get postgres connection: %w", err)
	}

	repo := a.providers.RepoProvider.GetRepository(db)
	storageClient := a.providers.StorageV1ClientProvider.GetStorageV1Client()

	uploader := a.providers.UploaderProvider.GetGatewayUploaderService(storageClient, repo, cfg.Concurrency, zapLogger)
	downloader := a.providers.DownloaderProvider.GetGatewayDownloaderService(storageClient, zapLogger)
	deleter := a.providers.DeleterProvider.GetGatewayDeleterService(storageClient, repo, cfg.Concurrency, zapLogger)

	service := a.providers.ServiceProvider.GetGatewayService(
		repo,
		uploader,
		downloader,
		deleter,
		cfg.NumChunks,
		cfg.MinFileSizeToSplit,
		zapLogger,
	)

	handler := a.providers.HandlerProvider.GetGatewayHandler(service, zapLogger)
	handler.AssignHandlers(a.echo)

	a.listenPort = cfg.HTTP.ListenPort

	return nil
}

// Run starts gateway application.
func (a *gatewayApp) Run(ctx context.Context) error {
	return runEchoWithGracefulShutdown(ctx, a.echo, a.listenPort, a.logger)
}
