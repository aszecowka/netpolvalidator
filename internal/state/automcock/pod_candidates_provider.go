// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/aszecowka/netpolvalidator/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// PodCandidatesProvider is an autogenerated mock type for the PodCandidatesProvider type
type PodCandidatesProvider struct {
	mock.Mock
}

// GetPodCandidatesForNamespace provides a mock function with given fields: ctx, ns
func (_m *PodCandidatesProvider) GetPodCandidatesForNamespace(ctx context.Context, ns string) ([]model.PodCandidate, error) {
	ret := _m.Called(ctx, ns)

	var r0 []model.PodCandidate
	if rf, ok := ret.Get(0).(func(context.Context, string) []model.PodCandidate); ok {
		r0 = rf(ctx, ns)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.PodCandidate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ns)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
