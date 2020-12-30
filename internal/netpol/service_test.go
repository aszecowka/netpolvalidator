package netpol_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes/fake"

	"github.com/aszecowka/netpolvalidator/internal/netpol"
)

func TestServiceGetNetworkPoliciesForNamespace(t *testing.T) {
	// GIVEN
	fakeClientset := fake.NewSimpleClientset(fixNetPolPaymentA(), fixNetPolPaymentB())
	// WHEN
	sut := netpol.NewService(fakeClientset.NetworkingV1())
	// THEN
	actual, err := sut.GetNetworkPoliciesForNamespace(context.Background(), "payment")
	require.NoError(t, err)
	require.Len(t, actual, 2)
	assert.Contains(t, actual, *fixNetPolPaymentA())
	assert.Contains(t, actual, *fixNetPolPaymentB())

}

func fixNetPolPaymentA() *v1.NetworkPolicy {
	return &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ingress-egress-for-payment-a",
			Namespace: "payment",
		},
	}
}

func fixNetPolPaymentB() *v1.NetworkPolicy {
	return &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ingress-egress-for-payment-b",
			Namespace: "payment",
		},
	}
}
