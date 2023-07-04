// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	echo "github.com/labstack/echo/v4"
	mock "github.com/stretchr/testify/mock"
)

// EchoProvider is an autogenerated mock type for the EchoProvider type
type EchoProvider struct {
	mock.Mock
}

// GetEcho provides a mock function with given fields:
func (_m *EchoProvider) GetEcho() *echo.Echo {
	ret := _m.Called()

	var r0 *echo.Echo
	if rf, ok := ret.Get(0).(func() *echo.Echo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*echo.Echo)
		}
	}

	return r0
}

type mockConstructorTestingTNewEchoProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewEchoProvider creates a new instance of EchoProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEchoProvider(t mockConstructorTestingTNewEchoProvider) *EchoProvider {
	mock := &EchoProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
