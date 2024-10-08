// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// UploadService is an autogenerated mock type for the UploadService type
type UploadService struct {
	mock.Mock
}

// UploadFile provides a mock function with given fields: file, filename, contentType
func (_m *UploadService) UploadFile(file io.Reader, filename string, contentType string) (string, error) {
	ret := _m.Called(file, filename, contentType)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(io.Reader, string, string) (string, error)); ok {
		return rf(file, filename, contentType)
	}
	if rf, ok := ret.Get(0).(func(io.Reader, string, string) string); ok {
		r0 = rf(file, filename, contentType)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(io.Reader, string, string) error); ok {
		r1 = rf(file, filename, contentType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUploadService creates a new instance of UploadService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUploadService(t interface {
	mock.TestingT
	Cleanup(func())
}) *UploadService {
	mock := &UploadService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
