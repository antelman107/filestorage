// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	domain "github.com/antelman107/filestorage/pkg/domain"
	mock "github.com/stretchr/testify/mock"

	zap "go.uber.org/zap"
)

// GatewayServiceProvider is an autogenerated mock type for the GatewayServiceProvider type
type GatewayServiceProvider struct {
	mock.Mock
}

// GetGatewayService provides a mock function with given fields: repo, uploader, downloader, deleter, numChunks, minFileSizeToSplit, logger
func (_m *GatewayServiceProvider) GetGatewayService(repo domain.GatewayRepository, uploader domain.GatewayUploaderService, downloader domain.GatewayDownloaderService, deleter domain.GatewayDeleterService, numChunks int64, minFileSizeToSplit int64, logger *zap.Logger) domain.GatewayService {
	ret := _m.Called(repo, uploader, downloader, deleter, numChunks, minFileSizeToSplit, logger)

	var r0 domain.GatewayService
	if rf, ok := ret.Get(0).(func(domain.GatewayRepository, domain.GatewayUploaderService, domain.GatewayDownloaderService, domain.GatewayDeleterService, int64, int64, *zap.Logger) domain.GatewayService); ok {
		r0 = rf(repo, uploader, downloader, deleter, numChunks, minFileSizeToSplit, logger)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.GatewayService)
		}
	}

	return r0
}

type mockConstructorTestingTNewGatewayServiceProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewGatewayServiceProvider creates a new instance of GatewayServiceProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGatewayServiceProvider(t mockConstructorTestingTNewGatewayServiceProvider) *GatewayServiceProvider {
	mock := &GatewayServiceProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
