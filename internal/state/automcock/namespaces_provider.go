// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	v1 "k8s.io/api/core/v1"
)

// NamespacesProvider is an autogenerated mock type for the NamespacesProvider type
type NamespacesProvider struct {
	mock.Mock
}

// GetAllNamespaces provides a mock function with given fields: ctx
func (_m *NamespacesProvider) GetAllNamespaces(ctx context.Context) ([]v1.Namespace, error) {
	ret := _m.Called(ctx)

	var r0 []v1.Namespace
	if rf, ok := ret.Get(0).(func(context.Context) []v1.Namespace); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]v1.Namespace)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}