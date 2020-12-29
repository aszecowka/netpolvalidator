package model

import (
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

type PodCandidate struct {
	OwnerName string
	Labels    map[string]string
}

type ClusterState struct {
	Namespaces      []v1.Namespace
	NetworkPolicies map[string][]networkingv1.NetworkPolicy
	PodCandidates   map[string][]PodCandidate
}
