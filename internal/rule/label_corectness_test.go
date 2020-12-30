package rule_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	v13 "k8s.io/api/networking/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/aszecowka/netpolvalidator/internal/model"
	"github.com/aszecowka/netpolvalidator/internal/rule"
)

const (
	labelDomain = "domain"
	labelApp    = "app"
	nsOrders    = "orders"
	nsPayments  = "payments"
)

func TestValidate(t *testing.T) {
	t.Run("no violations", func(t *testing.T) {
		// GIVEN
		sut := rule.NewLabelCorrectness()
		givenState := model.ClusterState{
			Namespaces: []v1.Namespace{fixNsOrders()},
			NetworkPolicies: map[string][]v13.NetworkPolicy{
				nsOrders: {
					fixIngressNetworkPolicyForOrdersA(),
				},
				nsPayments: {
					fixIngressNetworkPolicyForPaymentsA(),
				},
			},
			PodCandidates: map[string][]model.PodCandidate{
				nsOrders: {
					fixPodCandidateOrdersA(),
				},
				nsPayments: {
					fixPodCandidatePaymentsA(),
				},
			},
		}
		// WHEN
		violations, err := sut.Validate(givenState)
		// THEN
		require.NoError(t, err)
		require.Empty(t, violations)
	})

	t.Run("network policy pod selector does not match any pod", func(t *testing.T) {
		// GIVEN
		sut := rule.NewLabelCorrectness()
		givenState := model.ClusterState{
			Namespaces: []v1.Namespace{fixNsOrders()},
			NetworkPolicies: map[string][]v13.NetworkPolicy{
				nsOrders: {
					fixIngressNetworkPolicyForOrdersA(),
				},
			},
			PodCandidates: map[string][]model.PodCandidate{
				nsOrders: {
					fixPodCandidateOrdersB(),
				},
			},
		}
		// WHEN
		actual, err := sut.Validate(givenState)
		// THEN
		require.NoError(t, err)
		require.Len(t, actual, 1)
		assert.Equal(t, model.NewViolation(fixIngressNetworkPolicyForOrdersA(), "no pods matching pod selector", model.ViolationInvalidLabel), actual[0])
	})

	t.Run("returns many combined errors", func(t *testing.T) {
		// GIVEN
		sut := rule.NewLabelCorrectness()
		givenState := model.ClusterState{
			Namespaces: []v1.Namespace{fixNsOrders()},
			NetworkPolicies: map[string][]v13.NetworkPolicy{
				nsOrders: {
					fixIngressNetworkPolicyForOrdersA(),
				},
				nsPayments: {
					fixIngressNetworkPolicyForPaymentsA(),
				},
			},
			PodCandidates: map[string][]model.PodCandidate{
				nsOrders: {
					fixPodCandidateOrdersB(),
				},
				nsPayments: {
					fixPodCandidatePaymentsB(),
				},
			},
		}
		// WHEN
		actual, err := sut.Validate(givenState)
		// THEN
		require.NoError(t, err)
		require.NotEmpty(t, actual)
		// TODO more assertions
	})

}

func fixNsOrders() v1.Namespace {
	return v1.Namespace{
		ObjectMeta: v12.ObjectMeta{
			Name: nsOrders,
			Labels: map[string]string{
				labelDomain: nsOrders,
			},
		},
	}
}

func fixNsPayments() v1.Namespace {
	return v1.Namespace{
		ObjectMeta: v12.ObjectMeta{
			Name: nsPayments,
			Labels: map[string]string{
				labelDomain: nsPayments,
			},
		},
	}
}

func fixIngressNetworkPolicyForOrdersA() v13.NetworkPolicy {
	return v13.NetworkPolicy{
		ObjectMeta: v12.ObjectMeta{
			Name:      "ingress-for-orders-a",
			Namespace: nsOrders,
		},
		Spec: v13.NetworkPolicySpec{
			PodSelector: v12.LabelSelector{
				MatchLabels: map[string]string{
					labelApp: "orders-a",
				},
			},
		},
	}
}

func fixIngressNetworkPolicyForPaymentsA() v13.NetworkPolicy {
	return v13.NetworkPolicy{
		ObjectMeta: v12.ObjectMeta{
			Name:      "ingress-for-payments-a",
			Namespace: nsPayments,
		},
		Spec: v13.NetworkPolicySpec{
			PodSelector: v12.LabelSelector{
				MatchExpressions: []v12.LabelSelectorRequirement{
					{Key: labelApp, Operator: v12.LabelSelectorOpIn, Values: []string{"payments-a"}},
				},
			},
		},
	}
}

func fixPodCandidateOrdersA() model.PodCandidate {
	return model.PodCandidate{
		Labels: map[string]string{
			labelApp: "orders-a",
		},
	}
}

func fixPodCandidateOrdersB() model.PodCandidate {
	return model.PodCandidate{
		Labels: map[string]string{
			labelApp: "orders-b",
		},
	}
}

func fixPodCandidatePaymentsA() model.PodCandidate {
	return model.PodCandidate{
		Labels: map[string]string{
			labelApp: "payments-a",
		},
	}
}

func fixPodCandidatePaymentsB() model.PodCandidate {
	return model.PodCandidate{
		Labels: map[string]string{
			labelApp: "payments-b",
		},
	}
}
