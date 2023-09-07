// Code generated by mockery v2.23.2. DO NOT EDIT.

// Regenerate this file using `make filestore-mocks`.

package mocks

import (
	io "io"

	filestore "github.com/mattermost/mattermost/server/v8/platform/shared/filestore"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// FileBackend is an autogenerated mock type for the FileBackend type
type FileBackend struct {
	mock.Mock
}

// AppendFile provides a mock function with given fields: fr, path
func (_m *FileBackend) AppendFile(fr io.Reader, path string) (int64, error) {
	ret := _m.Called(fr, path)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(io.Reader, string) (int64, error)); ok {
		return rf(fr, path)
	}
	if rf, ok := ret.Get(0).(func(io.Reader, string) int64); ok {
		r0 = rf(fr, path)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(io.Reader, string) error); ok {
		r1 = rf(fr, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CopyFile provides a mock function with given fields: oldPath, newPath
func (_m *FileBackend) CopyFile(oldPath string, newPath string) error {
	ret := _m.Called(oldPath, newPath)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(oldPath, newPath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DriverName provides a mock function with given fields:
func (_m *FileBackend) DriverName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// FileExists provides a mock function with given fields: path
func (_m *FileBackend) FileExists(path string) (bool, error) {
	ret := _m.Called(path)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (bool, error)); ok {
		return rf(path)
	}
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FileModTime provides a mock function with given fields: path
func (_m *FileBackend) FileModTime(path string) (time.Time, error) {
	ret := _m.Called(path)

	var r0 time.Time
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (time.Time, error)); ok {
		return rf(path)
	}
	if rf, ok := ret.Get(0).(func(string) time.Time); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FileSize provides a mock function with given fields: path
func (_m *FileBackend) FileSize(path string) (int64, error) {
	ret := _m.Called(path)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (int64, error)); ok {
		return rf(path)
	}
	if rf, ok := ret.Get(0).(func(string) int64); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListDirectory provides a mock function with given fields: path
func (_m *FileBackend) ListDirectory(path string) ([]string, error) {
	ret := _m.Called(path)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]string, error)); ok {
		return rf(path)
	}
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListDirectoryRecursively provides a mock function with given fields: path
func (_m *FileBackend) ListDirectoryRecursively(path string) ([]string, error) {
	ret := _m.Called(path)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]string, error)); ok {
		return rf(path)
	}
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MoveFile provides a mock function with given fields: oldPath, newPath
func (_m *FileBackend) MoveFile(oldPath string, newPath string) error {
	ret := _m.Called(oldPath, newPath)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(oldPath, newPath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReadFile provides a mock function with given fields: path
func (_m *FileBackend) ReadFile(path string) ([]byte, error) {
	ret := _m.Called(path)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]byte, error)); ok {
		return rf(path)
	}
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Reader provides a mock function with given fields: path
func (_m *FileBackend) Reader(path string) (filestore.ReadCloseSeeker, error) {
	ret := _m.Called(path)

	var r0 filestore.ReadCloseSeeker
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (filestore.ReadCloseSeeker, error)); ok {
		return rf(path)
	}
	if rf, ok := ret.Get(0).(func(string) filestore.ReadCloseSeeker); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(filestore.ReadCloseSeeker)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveDirectory provides a mock function with given fields: path
func (_m *FileBackend) RemoveDirectory(path string) error {
	ret := _m.Called(path)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveFile provides a mock function with given fields: path
func (_m *FileBackend) RemoveFile(path string) error {
	ret := _m.Called(path)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TestConnection provides a mock function with given fields:
func (_m *FileBackend) TestConnection() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteFile provides a mock function with given fields: fr, path
func (_m *FileBackend) WriteFile(fr io.Reader, path string) (int64, error) {
	ret := _m.Called(fr, path)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(io.Reader, string) (int64, error)); ok {
		return rf(fr, path)
	}
	if rf, ok := ret.Get(0).(func(io.Reader, string) int64); ok {
		r0 = rf(fr, path)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(io.Reader, string) error); ok {
		r1 = rf(fr, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewFileBackend interface {
	mock.TestingT
	Cleanup(func())
}

// NewFileBackend creates a new instance of FileBackend. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFileBackend(t mockConstructorTestingTNewFileBackend) *FileBackend {
	mock := &FileBackend{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
