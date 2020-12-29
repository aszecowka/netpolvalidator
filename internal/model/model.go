package model

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

const (
	ViolationInvalidLabel ViolationType = "Invalid Label"
)

type ViolationType string

type PodCandidate struct {
	OwnerName string
	Labels    map[string]string
}

type ClusterState struct {
	Namespaces      []v1.Namespace
	NetworkPolicies map[string][]networkingv1.NetworkPolicy
	PodCandidates   map[string][]PodCandidate
}

type Violation struct {
	NetworkPolicyName      string
	NetworkPolicyNamespace string
	Message                string
	Type                   ViolationType
}

func NewViolation(ns, netPolName, message string, vType ViolationType) Violation {
	return Violation{
		NetworkPolicyNamespace: ns,
		NetworkPolicyName:      netPolName,
		Message:                message,
		Type:                   vType,
	}
}

func (v Violation) String() string {
	return fmt.Sprintf("[%s:%s]: %s: %s", v.NetworkPolicyNamespace, v.NetworkPolicyName, v.Type, v.Message)
}
