// Code generated by mockery v2.40.2. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/greenbone/opensight-notification-service/pkg/models"
	mock "github.com/stretchr/testify/mock"

	query "github.com/greenbone/opensight-golang-libraries/pkg/query"
)

// NotificationService is an autogenerated mock type for the NotificationService type
type NotificationService struct {
	mock.Mock
}

type NotificationService_Expecter struct {
	mock *mock.Mock
}

func (_m *NotificationService) EXPECT() *NotificationService_Expecter {
	return &NotificationService_Expecter{mock: &_m.Mock}
}

// CreateNotification provides a mock function with given fields: ctx, notificationIn
func (_m *NotificationService) CreateNotification(ctx context.Context, notificationIn models.Notification) (models.Notification, error) {
	ret := _m.Called(ctx, notificationIn)

	if len(ret) == 0 {
		panic("no return value specified for CreateNotification")
	}

	var r0 models.Notification
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Notification) (models.Notification, error)); ok {
		return rf(ctx, notificationIn)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.Notification) models.Notification); ok {
		r0 = rf(ctx, notificationIn)
	} else {
		r0 = ret.Get(0).(models.Notification)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.Notification) error); ok {
		r1 = rf(ctx, notificationIn)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NotificationService_CreateNotification_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateNotification'
type NotificationService_CreateNotification_Call struct {
	*mock.Call
}

// CreateNotification is a helper method to define mock.On call
//   - ctx context.Context
//   - notificationIn models.Notification
func (_e *NotificationService_Expecter) CreateNotification(ctx interface{}, notificationIn interface{}) *NotificationService_CreateNotification_Call {
	return &NotificationService_CreateNotification_Call{Call: _e.mock.On("CreateNotification", ctx, notificationIn)}
}

func (_c *NotificationService_CreateNotification_Call) Run(run func(ctx context.Context, notificationIn models.Notification)) *NotificationService_CreateNotification_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.Notification))
	})
	return _c
}

func (_c *NotificationService_CreateNotification_Call) Return(notification models.Notification, err error) *NotificationService_CreateNotification_Call {
	_c.Call.Return(notification, err)
	return _c
}

func (_c *NotificationService_CreateNotification_Call) RunAndReturn(run func(context.Context, models.Notification) (models.Notification, error)) *NotificationService_CreateNotification_Call {
	_c.Call.Return(run)
	return _c
}

// ListNotifications provides a mock function with given fields: ctx, resultSelector
func (_m *NotificationService) ListNotifications(ctx context.Context, resultSelector query.ResultSelector) ([]models.Notification, uint64, error) {
	ret := _m.Called(ctx, resultSelector)

	if len(ret) == 0 {
		panic("no return value specified for ListNotifications")
	}

	var r0 []models.Notification
	var r1 uint64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, query.ResultSelector) ([]models.Notification, uint64, error)); ok {
		return rf(ctx, resultSelector)
	}
	if rf, ok := ret.Get(0).(func(context.Context, query.ResultSelector) []models.Notification); ok {
		r0 = rf(ctx, resultSelector)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Notification)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, query.ResultSelector) uint64); ok {
		r1 = rf(ctx, resultSelector)
	} else {
		r1 = ret.Get(1).(uint64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, query.ResultSelector) error); ok {
		r2 = rf(ctx, resultSelector)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// NotificationService_ListNotifications_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListNotifications'
type NotificationService_ListNotifications_Call struct {
	*mock.Call
}

// ListNotifications is a helper method to define mock.On call
//   - ctx context.Context
//   - resultSelector query.ResultSelector
func (_e *NotificationService_Expecter) ListNotifications(ctx interface{}, resultSelector interface{}) *NotificationService_ListNotifications_Call {
	return &NotificationService_ListNotifications_Call{Call: _e.mock.On("ListNotifications", ctx, resultSelector)}
}

func (_c *NotificationService_ListNotifications_Call) Run(run func(ctx context.Context, resultSelector query.ResultSelector)) *NotificationService_ListNotifications_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(query.ResultSelector))
	})
	return _c
}

func (_c *NotificationService_ListNotifications_Call) Return(notifications []models.Notification, totalResult uint64, err error) *NotificationService_ListNotifications_Call {
	_c.Call.Return(notifications, totalResult, err)
	return _c
}

func (_c *NotificationService_ListNotifications_Call) RunAndReturn(run func(context.Context, query.ResultSelector) ([]models.Notification, uint64, error)) *NotificationService_ListNotifications_Call {
	_c.Call.Return(run)
	return _c
}

// NewNotificationService creates a new instance of NotificationService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewNotificationService(t interface {
	mock.TestingT
	Cleanup(func())
}) *NotificationService {
	mock := &NotificationService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}