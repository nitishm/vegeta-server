// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import dispatcher "vegeta-server/internal/dispatcher"
import io "io"
import mock "github.com/stretchr/testify/mock"

// ITaskActions is an autogenerated mock type for the ITaskActions type
type ITaskActions struct {
	mock.Mock
}

// Cancel provides a mock function with given fields:
func (_m *ITaskActions) Cancel() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Complete provides a mock function with given fields: _a0
func (_m *ITaskActions) Complete(_a0 io.Reader) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Reader) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Fail provides a mock function with given fields:
func (_m *ITaskActions) Fail() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Run provides a mock function with given fields: _a0
func (_m *ITaskActions) Run(_a0 dispatcher.AttackFunc) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(dispatcher.AttackFunc) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendUpdate provides a mock function with given fields:
func (_m *ITaskActions) SendUpdate() {
	_m.Called()
}
