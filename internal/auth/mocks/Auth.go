// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Auth is an autogenerated mock type for the Auth type
type Auth struct {
	mock.Mock
}

// Decode provides a mock function with given fields: sequence
func (_m *Auth) Decode(sequence []byte) (*string, error) {
	ret := _m.Called(sequence)

	var r0 *string
	if rf, ok := ret.Get(0).(func([]byte) *string); ok {
		r0 = rf(sequence)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(sequence)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Encode provides a mock function with given fields: id
func (_m *Auth) Encode(id *string) ([]byte, error) {
	ret := _m.Called(id)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(*string) []byte); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
