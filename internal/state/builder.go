package state

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"

	"github.com/aszecowka/netpolvalidator/internal/model"
)

//go:generate ${GOBIN}/mockery -name=NamespacesProvider -output=automcock -outpkg=automock -case=underscore
type NamespacesProvider interface {
	GetAllNamespaces(ctx context.Context) ([]v1.Namespace, error)
}

//go:generate ${GOBIN}/mockery -name=NetworkPoliciesProvider -output=automcock -outpkg=automock -case=underscore
type NetworkPoliciesProvider interface {
	GetNetworkPoliciesForNamespace(ctx context.Context, ns string) ([]netv1.NetworkPolicy, error)
}

//go:generate ${GOBIN}/mockery -name=PodCandidatesProvider -output=automcock -outpkg=automock -case=underscore
type PodCandidatesProvider interface {
	GetPodCandidatesForNamespace(ctx context.Context, ns string) ([]model.PodCandidate, error)
}

func NewBuilder(nsProvider NamespacesProvider, netPolProvider NetworkPoliciesProvider, podCandidatesProviders map[string]PodCandidatesProvider) *Builder {
	return &Builder{
		nsProvider:             nsProvider,
		netPolProvider:         netPolProvider,
		podCandidatesProviders: podCandidatesProviders,
	}
}

type Builder struct {
	nsProvider             NamespacesProvider
	netPolProvider         NetworkPoliciesProvider
	podCandidatesProviders map[string]PodCandidatesProvider
}

func (b *Builder) Build(ctx context.Context) (*model.ClusterState, error) {
	out := &model.ClusterState{}
	namespaces, err := b.nsProvider.GetAllNamespaces(ctx)
	if err != nil {
		return nil, fmt.Errorf("while getting all namespaces: %w", err)
	}
	out.Namespaces = namespaces
	out.NetworkPolicies = make(map[string][]netv1.NetworkPolicy)

	for _, ns := range namespaces {
		policies, err := b.netPolProvider.GetNetworkPoliciesForNamespace(ctx, ns.Name)
		if err != nil {
			return nil, fmt.Errorf("while getting network policies for namespace: %s: %w", ns.Name, err)
		}
		out.NetworkPolicies[ns.Name] = policies
	}

	out.PodCandidates = make(map[string][]model.PodCandidate)
	for podCandidateStrategyName, strategy := range b.podCandidatesProviders {
		for _, ns := range namespaces {
			podCandidates, err := strategy.GetPodCandidatesForNamespace(ctx, ns.Name)
			if err != nil {
				return nil, fmt.Errorf("while getting pod candidates for namespace: %s, strategy: %s: %w", ns.Name, podCandidateStrategyName, err)
			}
			if out.PodCandidates[ns.Name] == nil {
				out.PodCandidates[ns.Name] = podCandidates
			} else {
				out.PodCandidates[ns.Name] = append(out.PodCandidates[ns.Name], podCandidates...)
			}
		}
	}
	return out, nil
}
