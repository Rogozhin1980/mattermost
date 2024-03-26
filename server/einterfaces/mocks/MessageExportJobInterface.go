// Code generated by mockery v2.23.2. DO NOT EDIT.

// Regenerate this file using `make einterfaces-mocks`.

package mocks

import (
	jobs "github.com/mattermost/mattermost/server/v8/einterfaces/jobs"
	mock "github.com/stretchr/testify/mock"

	model "github.com/mattermost/mattermost/server/public/model"
)

// MessageExportJobInterface is an autogenerated mock type for the MessageExportJobInterface type
type MessageExportJobInterface struct {
	mock.Mock
}

// MakeScheduler provides a mock function with given fields:
func (_m *MessageExportJobInterface) MakeScheduler() jobs.Scheduler {
	ret := _m.Called()

	var r0 jobs.Scheduler
	if rf, ok := ret.Get(0).(func() jobs.Scheduler); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(jobs.Scheduler)
		}
	}

	return r0
}

// MakeWorker provides a mock function with given fields:
func (_m *MessageExportJobInterface) MakeWorker() model.Worker {
	ret := _m.Called()

	var r0 model.Worker
	if rf, ok := ret.Get(0).(func() model.Worker); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.Worker)
		}
	}

	return r0
}

type mockConstructorTestingTNewMessageExportJobInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewMessageExportJobInterface creates a new instance of MessageExportJobInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMessageExportJobInterface(t mockConstructorTestingTNewMessageExportJobInterface) *MessageExportJobInterface {
	mock := &MessageExportJobInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
