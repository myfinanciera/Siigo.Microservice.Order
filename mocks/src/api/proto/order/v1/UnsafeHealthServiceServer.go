// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UnsafeHealthServiceServer is an autogenerated mock type for the UnsafeHealthServiceServer type
type UnsafeHealthServiceServer struct {
	mock.Mock
}

// mustEmbedUnimplementedHealthServiceServer provides a mock function with given fields:
func (_m *UnsafeHealthServiceServer) mustEmbedUnimplementedHealthServiceServer() {
	_m.Called()
}

type mockConstructorTestingTNewUnsafeHealthServiceServer interface {
	mock.TestingT
	Cleanup(func())
}

// NewUnsafeHealthServiceServer creates a new instance of UnsafeHealthServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUnsafeHealthServiceServer(t mockConstructorTestingTNewUnsafeHealthServiceServer) *UnsafeHealthServiceServer {
	mock := &UnsafeHealthServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}