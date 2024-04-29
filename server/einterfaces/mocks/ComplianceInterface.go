// Code generated by mockery v2.42.2. DO NOT EDIT.

// Regenerate this file using `make einterfaces-mocks`.

package mocks

import (
	model "github.com/mattermost/mattermost/server/public/model"
	request "github.com/mattermost/mattermost/server/public/shared/request"
	mock "github.com/stretchr/testify/mock"
)

// ComplianceInterface is an autogenerated mock type for the ComplianceInterface type
type ComplianceInterface struct {
	mock.Mock
}

// RunComplianceJob provides a mock function with given fields: rctx, job
func (_m *ComplianceInterface) RunComplianceJob(rctx request.CTX, job *model.Compliance) *model.AppError {
	ret := _m.Called(rctx, job)

	if len(ret) == 0 {
		panic("no return value specified for RunComplianceJob")
	}

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(request.CTX, *model.Compliance) *model.AppError); ok {
		r0 = rf(rctx, job)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// StartComplianceDailyJob provides a mock function with given fields:
func (_m *ComplianceInterface) StartComplianceDailyJob() {
	_m.Called()
}

// NewComplianceInterface creates a new instance of ComplianceInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewComplianceInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *ComplianceInterface {
	mock := &ComplianceInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
