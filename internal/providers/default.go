package providers

import "github.com/antelman107/filestorage/internal/services"

var DefaultGatewayAppProviders = services.GatewayAppProviders{
	EchoProvider:            NewDefaultEchoProvider(),
	ConfigLoader:            NewDefaultConfigLoader(),
	SqlxProvider:            NewDefaultSQLXProvider(),
	RepoProvider:            NewDefaultGatewayRepositoryProvider(),
	StorageV1ClientProvider: NewDefaultStorageClientProvider(),
	UploaderProvider:        NewDefaultGatewayUploaderProvider(),
	DownloaderProvider:      NewDefaultGatewayDownloaderProvider(),
	DeleterProvider:         NewDefaultGatewayDeleterProvider(),
	ServiceProvider:         NewDefaultGatewayServiceProvider(),
	HandlerProvider:         NewDefaultGatewayHandlerProvider(),
}

var DefaultStorageAppProviders = services.StorageAppProviders{
	EchoProvider:    NewDefaultEchoProvider(),
	ConfigLoader:    NewDefaultConfigLoader(),
	ServiceProvider: NewDefaultStorageServiceProvider(),
	HandlerProvider: NewDefaultStorageHandlerProvider(),
}
