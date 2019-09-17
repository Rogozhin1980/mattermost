// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	gorp "github.com/mattermost/gorp"
	mock "github.com/stretchr/testify/mock"
)

// SqlSupplier is an autogenerated mock type for the SqlSupplier type
type SqlSupplier struct {
	mock.Mock
}

// GetMaster provides a mock function with given fields:
func (_m *SqlSupplier) GetMaster() *gorp.DbMap {
	ret := _m.Called()

	var r0 *gorp.DbMap
	if rf, ok := ret.Get(0).(func() *gorp.DbMap); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorp.DbMap)
		}
	}

	return r0
}
