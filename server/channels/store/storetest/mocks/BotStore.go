// Code generated by mockery v2.42.2. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	model "github.com/mattermost/mattermost/server/public/model"
	mock "github.com/stretchr/testify/mock"
)

// BotStore is an autogenerated mock type for the BotStore type
type BotStore struct {
	mock.Mock
}

// Get provides a mock function with given fields: userID, includeDeleted
func (_m *BotStore) Get(userID string, includeDeleted bool) (*model.Bot, error) {
	ret := _m.Called(userID, includeDeleted)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *model.Bot
	var r1 error
	if rf, ok := ret.Get(0).(func(string, bool) (*model.Bot, error)); ok {
		return rf(userID, includeDeleted)
	}
	if rf, ok := ret.Get(0).(func(string, bool) *model.Bot); ok {
		r0 = rf(userID, includeDeleted)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Bot)
		}
	}

	if rf, ok := ret.Get(1).(func(string, bool) error); ok {
		r1 = rf(userID, includeDeleted)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: options
func (_m *BotStore) GetAll(options *model.BotGetOptions) ([]*model.Bot, error) {
	ret := _m.Called(options)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []*model.Bot
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.BotGetOptions) ([]*model.Bot, error)); ok {
		return rf(options)
	}
	if rf, ok := ret.Get(0).(func(*model.BotGetOptions) []*model.Bot); ok {
		r0 = rf(options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Bot)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.BotGetOptions) error); ok {
		r1 = rf(options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllAfter provides a mock function with given fields: limit, afterId
func (_m *BotStore) GetAllAfter(limit int, afterId string) ([]*model.Bot, error) {
	ret := _m.Called(limit, afterId)

	if len(ret) == 0 {
		panic("no return value specified for GetAllAfter")
	}

	var r0 []*model.Bot
	var r1 error
	if rf, ok := ret.Get(0).(func(int, string) ([]*model.Bot, error)); ok {
		return rf(limit, afterId)
	}
	if rf, ok := ret.Get(0).(func(int, string) []*model.Bot); ok {
		r0 = rf(limit, afterId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Bot)
		}
	}

	if rf, ok := ret.Get(1).(func(int, string) error); ok {
		r1 = rf(limit, afterId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByUsername provides a mock function with given fields: username
func (_m *BotStore) GetByUsername(username string) (*model.Bot, error) {
	ret := _m.Called(username)

	if len(ret) == 0 {
		panic("no return value specified for GetByUsername")
	}

	var r0 *model.Bot
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.Bot, error)); ok {
		return rf(username)
	}
	if rf, ok := ret.Get(0).(func(string) *model.Bot); ok {
		r0 = rf(username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Bot)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PermanentDelete provides a mock function with given fields: userID
func (_m *BotStore) PermanentDelete(userID string) error {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for PermanentDelete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Save provides a mock function with given fields: bot
func (_m *BotStore) Save(bot *model.Bot) (*model.Bot, error) {
	ret := _m.Called(bot)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 *model.Bot
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.Bot) (*model.Bot, error)); ok {
		return rf(bot)
	}
	if rf, ok := ret.Get(0).(func(*model.Bot) *model.Bot); ok {
		r0 = rf(bot)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Bot)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.Bot) error); ok {
		r1 = rf(bot)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: bot
func (_m *BotStore) Update(bot *model.Bot) (*model.Bot, error) {
	ret := _m.Called(bot)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *model.Bot
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.Bot) (*model.Bot, error)); ok {
		return rf(bot)
	}
	if rf, ok := ret.Get(0).(func(*model.Bot) *model.Bot); ok {
		r0 = rf(bot)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Bot)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.Bot) error); ok {
		r1 = rf(bot)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewBotStore creates a new instance of BotStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBotStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *BotStore {
	mock := &BotStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
