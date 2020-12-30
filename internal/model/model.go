package model

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

const (
	ViolationInvalidLabel ViolationType = "Invalid Label"
	Ingress               RuleType      = "Ingress"
	Egress                RuleType      = "Egress"
)

type ViolationType string
type RuleType string

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

func NewViolation(np networkingv1.NetworkPolicy, message string, vType ViolationType) Violation {
	return Violation{
		NetworkPolicyNamespace: np.Namespace,
		NetworkPolicyName:      np.Name,
		Message:                message,
		Type:                   vType,
	}
}

func (v Violation) String() string {
	return fmt.Sprintf("[%s:%s]: %s: %s", v.NetworkPolicyNamespace, v.NetworkPolicyName, v.Type, v.Message)
}
