package state_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	v13 "k8s.io/api/networking/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/aszecowka/netpolvalidator/internal/model"
	"github.com/aszecowka/netpolvalidator/internal/state"
	automock "github.com/aszecowka/netpolvalidator/internal/state/automcock"
)

func TestBuildClusterState(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// GIVEN
		mockNsProvider := &automock.NamespacesProvider{}
		defer mockNsProvider.AssertExpectations(t)

		mockNsProvider.On("GetAllNamespaces", mock.Anything).Return([]v1.Namespace{fixNsWithName("a"), fixNsWithName("b")}, nil).Once()

		mockNetPolProvider := &automock.NetworkPoliciesProvider{}
		defer mockNetPolProvider.AssertExpectations(t)
		mockNetPolProvider.On("GetNetworkPoliciesForNamespace", mock.Anything, "a").Return([]v13.NetworkPolicy{fixFixNetPol("a", "ingress-deny-all"), fixFixNetPol("a", "egress-deny-all")}, nil).Once()
		mockNetPolProvider.On("GetNetworkPoliciesForNamespace", mock.Anything, "b").Return(nil, nil).Once()
		mockDeploymentsProvider := &automock.PodCandidatesProvider{}
		defer mockDeploymentsProvider.AssertExpectations(t)
		mockDeploymentsProvider.On("GetPodCandidatesForNamespace", mock.Anything, "a").Return([]model.PodCandidate{fixPodCandidate("deploy-a")}, nil).Once()
		mockDeploymentsProvider.On("GetPodCandidatesForNamespace", mock.Anything, "b").Return([]model.PodCandidate{fixPodCandidate("deploy-b")}, nil).Once()

		mockCronjobProvider := &automock.PodCandidatesProvider{}
		defer mockCronjobProvider.AssertExpectations(t)
		mockCronjobProvider.On("GetPodCandidatesForNamespace", mock.Anything, "a").Return([]model.PodCandidate{fixPodCandidate("cronjob-a")}, nil).Once()
		mockCronjobProvider.On("GetPodCandidatesForNamespace", mock.Anything, "b").Return(nil, nil).Once()

		podCandidatesProviders := map[string]state.PodCandidatesProvider{
			"deploy":  mockDeploymentsProvider,
			"cronjob": mockCronjobProvider,
		}

		sut := state.NewBuilder(mockNsProvider, mockNetPolProvider, podCandidatesProviders)
		// WHEN
		actual, err := sut.Build(context.Background())
		// THEN
		require.NoError(t, err)
		require.NotNil(t, actual)
		assert.Equal(t, []v1.Namespace{fixNsWithName("a"), fixNsWithName("b")}, actual.Namespaces)
		assert.Len(t, actual.NetworkPolicies, 2)
		assert.Equal(t, []v13.NetworkPolicy{fixFixNetPol("a", "ingress-deny-all"), fixFixNetPol("a", "egress-deny-all")}, actual.NetworkPolicies["a"])
		assert.Equal(t, []v13.NetworkPolicy(nil), actual.NetworkPolicies["b"])
		assert.Len(t, actual.PodCandidates, 2)
		assert.Len(t, actual.PodCandidates["a"], 2)
		assert.Contains(t, actual.PodCandidates["a"], fixPodCandidate("deploy-a"))
		assert.Contains(t, actual.PodCandidates["a"], fixPodCandidate("cronjob-a"))
		assert.Len(t, actual.PodCandidates["b"], 1)
		assert.Contains(t, actual.PodCandidates["b"], fixPodCandidate("deploy-b"))
	})

	t.Run("got error on getting namespaces", func(t *testing.T) {
		// GIVEN
		mockNsProvider := &automock.NamespacesProvider{}
		defer mockNsProvider.AssertExpectations(t)

		mockNsProvider.On("GetAllNamespaces", mock.Anything).Return(nil, errors.New("some error")).Once()
		sut := state.NewBuilder(mockNsProvider, nil, nil)
		// WHEN
		_, err := sut.Build(context.Background())
		// THEN
		require.EqualError(t, err, "while getting all namespaces: some error")
	})

	t.Run("got error on getting network policies", func(t *testing.T) {
		// GIVEN
		mockNsProvider := &automock.NamespacesProvider{}
		defer mockNsProvider.AssertExpectations(t)

		mockNsProvider.On("GetAllNamespaces", mock.Anything).Return([]v1.Namespace{fixNsWithName("a"), fixNsWithName("b")}, nil).Once()

		mockNetPolProvider := &automock.NetworkPoliciesProvider{}
		defer mockNetPolProvider.AssertExpectations(t)
		mockNetPolProvider.On("GetNetworkPoliciesForNamespace", mock.Anything, "a").Return(nil, errors.New("some error")).Once()
		sut := state.NewBuilder(mockNsProvider, mockNetPolProvider, nil)
		// WHEN
		_, err := sut.Build(context.Background())
		// THEN
		require.EqualError(t, err, "while getting network policies for namespace: a: some error")

	})

	t.Run("got error on getting pod candidates", func(t *testing.T) {
		// GIVEN
		mockNsProvider := &automock.NamespacesProvider{}
		defer mockNsProvider.AssertExpectations(t)
		mockNsProvider.On("GetAllNamespaces", mock.Anything).Return([]v1.Namespace{fixNsWithName("a")}, nil).Once()

		mockNetPolProvider := &automock.NetworkPoliciesProvider{}
		defer mockNetPolProvider.AssertExpectations(t)
		mockNetPolProvider.On("GetNetworkPoliciesForNamespace", mock.Anything, "a").Return(nil, nil).Once()
		mockDeploymentsProvider := &automock.PodCandidatesProvider{}
		defer mockDeploymentsProvider.AssertExpectations(t)
		mockDeploymentsProvider.On("GetPodCandidatesForNamespace", mock.Anything, "a").Return(nil, errors.New("some error")).Once()
		providers := map[string]state.PodCandidatesProvider{
			"deploy": mockDeploymentsProvider,
		}
		sut := state.NewBuilder(mockNsProvider, mockNetPolProvider, providers)
		// WHEN
		_, err := sut.Build(context.Background())
		require.EqualError(t, err, "while getting pod candidates for namespace: a, strategy: deploy: some error")
	})
}

func fixNsWithName(name string) v1.Namespace {
	return v1.Namespace{
		ObjectMeta: v12.ObjectMeta{
			Name: name,
		},
	}
}

func fixFixNetPol(namespace, name string) v13.NetworkPolicy {
	return v13.NetworkPolicy{
		ObjectMeta: v12.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func fixPodCandidate(ownerName string) model.PodCandidate {
	return model.PodCandidate{
		OwnerName: ownerName,
	}
}
