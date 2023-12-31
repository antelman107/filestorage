// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"
	http "net/http"

	domain "github.com/antelman107/filestorage/pkg/domain"

	mock "github.com/stretchr/testify/mock"
)

// StorageV1Client is an autogenerated mock type for the StorageV1Client type
type StorageV1Client struct {
	mock.Mock
}

// DeleteChunk provides a mock function with given fields: ctx, chunk
func (_m *StorageV1Client) DeleteChunk(ctx context.Context, chunk domain.Chunk) error {
	ret := _m.Called(ctx, chunk)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Chunk) error); ok {
		r0 = rf(ctx, chunk)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetChunk provides a mock function with given fields: ctx, chunk
func (_m *StorageV1Client) GetChunk(ctx context.Context, chunk domain.Chunk) (*http.Response, error) {
	ret := _m.Called(ctx, chunk)

	var r0 *http.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Chunk) (*http.Response, error)); ok {
		return rf(ctx, chunk)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.Chunk) *http.Response); ok {
		r0 = rf(ctx, chunk)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.Chunk) error); ok {
		r1 = rf(ctx, chunk)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHealth provides a mock function with given fields: ctx, serverURL
func (_m *StorageV1Client) GetHealth(ctx context.Context, serverURL string) error {
	ret := _m.Called(ctx, serverURL)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, serverURL)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PostChunk provides a mock function with given fields: ctx, chunk
func (_m *StorageV1Client) PostChunk(ctx context.Context, chunk domain.ChunkWithData) error {
	ret := _m.Called(ctx, chunk)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.ChunkWithData) error); ok {
		r0 = rf(ctx, chunk)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewStorageV1Client interface {
	mock.TestingT
	Cleanup(func())
}

// NewStorageV1Client creates a new instance of StorageV1Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStorageV1Client(t mockConstructorTestingTNewStorageV1Client) *StorageV1Client {
	mock := &StorageV1Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
