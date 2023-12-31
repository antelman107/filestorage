// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	sqlx "github.com/jmoiron/sqlx"
	mock "github.com/stretchr/testify/mock"
)

// SQLXProvider is an autogenerated mock type for the SQLXProvider type
type SQLXProvider struct {
	mock.Mock
}

// GetSQLX provides a mock function with given fields: connectionString
func (_m *SQLXProvider) GetSQLX(connectionString string) (*sqlx.DB, error) {
	ret := _m.Called(connectionString)

	var r0 *sqlx.DB
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*sqlx.DB, error)); ok {
		return rf(connectionString)
	}
	if rf, ok := ret.Get(0).(func(string) *sqlx.DB); ok {
		r0 = rf(connectionString)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqlx.DB)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(connectionString)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewSQLXProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewSQLXProvider creates a new instance of SQLXProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSQLXProvider(t mockConstructorTestingTNewSQLXProvider) *SQLXProvider {
	mock := &SQLXProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
