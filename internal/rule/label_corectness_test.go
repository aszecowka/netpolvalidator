package rule_test

import (
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
	sut := rule.NewLabelCorrectness()
	t.Run("pod selector is correct", func(t *testing.T) {
		// GIVEN
		givenState := model.ClusterState{
			Namespaces: []v1.Namespace{fixNsOrders()},
			NetworkPolicies: map[string][]netv1.NetworkPolicy{
				nsOrders: {fixIngressNetworkPolicyForOrdersA()},
			},
			PodCandidates: map[string][]model.PodCandidate{
				nsOrders: {
					fixPodCandidateOrdersA(),
				},
			},
		}
		// WHEN
		violations, err := sut.Validate(givenState)
		// THEN
		require.NoError(t, err)
		require.Empty(t, violations)
	})

	t.Run("pod selector does not match any pod", func(t *testing.T) {
		// GIVEN
		givenState := model.ClusterState{
			Namespaces: []v1.Namespace{fixNsOrders()},
			NetworkPolicies: map[string][]netv1.NetworkPolicy{
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

	t.Run("ingress rule for specific pods and namespaces is correct", func(t *testing.T) {
		// GIVEN
		givenState := model.ClusterState{
			Namespaces: []v1.Namespace{fixNsPayments(), fixNsOrders()},
			NetworkPolicies: map[string][]netv1.NetworkPolicy{
				nsPayments: {
					getNetPol(t, `
metadata:
  name: ingress-for-payments-a
  namespace: payments
spec:
  podSelector:
    matchLabels:
      app: payments-a
  ingress:
    - from:
      - namespaceSelector:
          matchLabels:
            domain: orders
        podSelector:
          matchLabels:
            app: orders-a
    `),
				},
			},
			PodCandidates: map[string][]model.PodCandidate{
				nsPayments: {fixPodCandidatePaymentsA()},
				nsOrders:   {fixPodCandidateOrdersA()},
			},
		}

		// WHEN
		actualViolations, err := sut.Validate(givenState)
		// THEN
		require.NoError(t, err)
		require.Empty(t, actualViolations)

	})

	t.Run("ingress rule for specific pods and namespaces does not match any namespace", func(t *testing.T) {
		// GIVEN
		givenNetPol := getNetPol(t, `
metadata:
  name: ingress-for-payments-a
  namespace: payments
spec:
  podSelector:
    matchLabels:
      app: payments-a
  ingress:
    - from:
      - namespaceSelector:
          matchLabels:
            domain: doesnotexist
        podSelector:
          matchLabels:
            app: orders-a
    `)
		givenState := model.ClusterState{
			Namespaces: []v1.Namespace{fixNsPayments(), fixNsOrders()},
			NetworkPolicies: map[string][]netv1.NetworkPolicy{
				nsPayments: {givenNetPol},
			},
			PodCandidates: map[string][]model.PodCandidate{
				nsPayments: {fixPodCandidatePaymentsA()},
				nsOrders:   {fixPodCandidateOrdersA()},
			},
		}

		// WHEN
		actualViolations, err := sut.Validate(givenState)
		require.NoError(t, err)
		require.Len(t, actualViolations, 1)
		assert.Equal(t, model.NewViolation(givenNetPol, "no namespaces matching labels for ingress rule [1:1]", model.ViolationInvalidLabel), actualViolations[0])
	})

	t.Run("ingress rule for specific pods and namespaces does not match any pods", func(t *testing.T) {
		// GIVEN
		givenNetPol := getNetPol(t, `
metadata:
  name: ingress-for-payments-a
  namespace: payments
spec:
  podSelector:
    matchLabels:
      app: payments-a
  ingress:
    - from:
      - namespaceSelector:
          matchLabels:
            domain: orders
        podSelector:
          matchLabels:
            app: orders-a
    `)
		givenState := model.ClusterState{
			Namespaces: []v1.Namespace{fixNsPayments(), fixNsOrders()},
			NetworkPolicies: map[string][]netv1.NetworkPolicy{
				nsPayments: {givenNetPol},
			},
			PodCandidates: map[string][]model.PodCandidate{
				nsPayments: {fixPodCandidatePaymentsA()},
				nsOrders:   {fixPodCandidateOrdersB()},
			},
		}

		// WHEN
		actualViolations, err := sut.Validate(givenState)
		// THEN
		require.NoError(t, err)
		require.Len(t, actualViolations, 1)
		assert.Equal(t, model.NewViolation(givenNetPol, "no pods matching labels for ingress rule [1:1]", model.ViolationInvalidLabel), actualViolations[0])
	})

	t.Run("ingress rule for pods in the network policy namespace is correct", func(t *testing.T) {

	})

	t.Run("ingress rule for pods in the network policy namespace does not match any pods", func(t *testing.T) {

	})

	t.Run("ingress rule for all pods in the selected namespaces is correct", func(t *testing.T) {

	})

	t.Run("ingress rule for all pods in the selected namespaces does not match any namespaces", func(t *testing.T) {

	})

	t.Run("ingress rule for all pods in the selected namespaces does not match any pod", func(t *testing.T) {

	})

	t.Run("returns many combined errors", func(t *testing.T) {
		// GIVEN
		givenState := model.ClusterState{
			Namespaces: []v1.Namespace{fixNsOrders()},
			NetworkPolicies: map[string][]netv1.NetworkPolicy{
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
		ObjectMeta: metav1.ObjectMeta{
			Name: nsOrders,
			Labels: map[string]string{
				labelDomain: nsOrders,
			},
		},
	}
}

func fixNsPayments() v1.Namespace {
	return v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: nsPayments,
			Labels: map[string]string{
				labelDomain: nsPayments,
			},
		},
	}
}

func fixIngressNetworkPolicyForOrdersA() netv1.NetworkPolicy {
	return netv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ingress-for-orders-a",
			Namespace: nsOrders,
		},
		Spec: netv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					labelApp: "orders-a",
				},
			},
		},
	}
}

func fixIngressNetworkPolicyForPaymentsA() netv1.NetworkPolicy {
	return netv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ingress-for-payments-a",
			Namespace: nsPayments,
		},
		Spec: netv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{Key: labelApp, Operator: metav1.LabelSelectorOpIn, Values: []string{"payments-a"}},
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

func getNetPol(t *testing.T, in string) netv1.NetworkPolicy {
	np := netv1.NetworkPolicy{}
	err := yaml.Unmarshal([]byte(in), &np)
	require.NoError(t, err)
	return np
}
