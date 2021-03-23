// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make einterfaces-mocks`.

package mocks

import (
	model "github.com/mattermost/mattermost-server/v5/model"
	mock "github.com/stretchr/testify/mock"
)

// SamlInterface is an autogenerated mock type for the SamlInterface type
type SamlInterface struct {
	mock.Mock
}

// BuildRequest provides a mock function with given fields: relayState
func (_m *SamlInterface) BuildRequest(relayState string) (*model.SamlAuthRequest, *model.AppError) {
	ret := _m.Called(relayState)

	var r0 *model.SamlAuthRequest
	if rf, ok := ret.Get(0).(func(string) *model.SamlAuthRequest); ok {
		r0 = rf(relayState)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.SamlAuthRequest)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(relayState)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// CheckProviderAttributes provides a mock function with given fields: SS, ouser, patch
func (_m *SamlInterface) CheckProviderAttributes(SS *model.SamlSettings, ouser *model.User, patch *model.UserPatch) string {
	ret := _m.Called(SS, ouser, patch)

	var r0 string
	if rf, ok := ret.Get(0).(func(*model.SamlSettings, *model.User, *model.UserPatch) string); ok {
		r0 = rf(SS, ouser, patch)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ConfigureSP provides a mock function with given fields:
func (_m *SamlInterface) ConfigureSP() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DoLogin provides a mock function with given fields: encodedXML, relayState
func (_m *SamlInterface) DoLogin(encodedXML string, relayState map[string]string) (*model.User, *model.AppError) {
	ret := _m.Called(encodedXML, relayState)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(string, map[string]string) *model.User); ok {
		r0 = rf(encodedXML, relayState)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, map[string]string) *model.AppError); ok {
		r1 = rf(encodedXML, relayState)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetMetadata provides a mock function with given fields:
func (_m *SamlInterface) GetMetadata() (string, *model.AppError) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func() *model.AppError); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// ResetAuthDataToEmail provides a mock function with given fields: includeDeleted, dryRun, specifiedUserIDs
func (_m *SamlInterface) ResetAuthDataToEmail(includeDeleted bool, dryRun bool, specifiedUserIDs []string) (int64, error) {
	ret := _m.Called(includeDeleted, dryRun, specifiedUserIDs)

	var r0 int64
	if rf, ok := ret.Get(0).(func(bool, bool, []string) int64); ok {
		r0 = rf(includeDeleted, dryRun, specifiedUserIDs)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool, bool, []string) error); ok {
		r1 = rf(includeDeleted, dryRun, specifiedUserIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
