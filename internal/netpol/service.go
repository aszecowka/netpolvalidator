package netpol

import (
	"context"
	"fmt"

	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typednetv1 "k8s.io/client-go/kubernetes/typed/networking/v1"
)

type service struct {
	client typednetv1.NetworkPoliciesGetter
}

func NewService(client typednetv1.NetworkPoliciesGetter) *service {
	return &service{client: client}
}

func (s *service) GetNetworkPoliciesForNamespace(ctx context.Context, ns string) ([]netv1.NetworkPolicy, error) {
	list, err := s.client.NetworkPolicies(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("while listing network policies from namespace: %s: %w", ns, err)
	}
	// TODO continue
	return list.Items, nil
}
