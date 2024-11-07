// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// IPList is an autogenerated mock type for the IPList type
type IPList struct {
	mock.Mock
}

// BlacklistCheckIP provides a mock function with given fields: ctx, ip
func (_m *IPList) BlacklistCheckIP(ctx context.Context, ip string) (bool, error) {
	ret := _m.Called(ctx, ip)

	if len(ret) == 0 {
		panic("no return value specified for BlacklistCheckIP")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (bool, error)); ok {
		return rf(ctx, ip)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, ip)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ip)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WhitelistCheckIP provides a mock function with given fields: ctx, ip
func (_m *IPList) WhitelistCheckIP(ctx context.Context, ip string) (bool, error) {
	ret := _m.Called(ctx, ip)

	if len(ret) == 0 {
		panic("no return value specified for WhitelistCheckIP")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (bool, error)); ok {
		return rf(ctx, ip)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, ip)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ip)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIPList creates a new instance of IPList. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIPList(t interface {
	mock.TestingT
	Cleanup(func())
}) *IPList {
	mock := &IPList{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
