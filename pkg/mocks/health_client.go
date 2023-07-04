// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// HealthClient is an autogenerated mock type for the HealthClient type
type HealthClient struct {
	mock.Mock
}

// GetHealth provides a mock function with given fields: ctx, serverURL
func (_m *HealthClient) GetHealth(ctx context.Context, serverURL string) error {
	ret := _m.Called(ctx, serverURL)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, serverURL)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewHealthClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewHealthClient creates a new instance of HealthClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewHealthClient(t mockConstructorTestingTNewHealthClient) *HealthClient {
	mock := &HealthClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
