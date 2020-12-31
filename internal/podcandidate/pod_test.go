package podcandidate_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/aszecowka/netpolvalidator/internal/model"
	"github.com/aszecowka/netpolvalidator/internal/podcandidate"
)

func TestPodCandidateFromPods(t *testing.T) {
	// GIVEN
	fakeClientset := fake.NewSimpleClientset(fixPodA(), fixPodB())
	sut := podcandidate.NewPodsFetcher(fakeClientset.CoreV1())
	// WHEN
	actual, err := sut.GetPodCandidatesForNamespace(context.Background(), "orders")
	// THEN
	require.NoError(t, err)
	require.Len(t, actual, 2)
	assert.Contains(t, actual, model.PodCandidate{OwnerName: "pod/orders/pod-a", Labels: map[string]string{
		"app": "app-a",
	}})
	assert.Contains(t, actual, model.PodCandidate{OwnerName: "pod/orders/pod-b", Labels: map[string]string{
		"app": "app-b",
	}})
}

func fixPodA() *v1.Pod {
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod-a",
			Namespace: "orders",
			Labels: map[string]string{
				"app": "app-a",
			},
		},
	}
}

func fixPodB() *v1.Pod {
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod-b",
			Namespace: "orders",
			Labels: map[string]string{
				"app": "app-b",
			},
		},
	}
}
