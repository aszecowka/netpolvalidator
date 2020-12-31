package podcandidate_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/aszecowka/netpolvalidator/internal/model"
	"github.com/aszecowka/netpolvalidator/internal/podcandidate"
)

func TestPodCandidateFromStatefulsets(t *testing.T) {
	// GIVEN
	fakeClientset := fake.NewSimpleClientset(fixStatefulsetA(), fixStatefulsetB())
	sut := podcandidate.NewStatefulsetsFetcher(fakeClientset.AppsV1())
	// WHEN
	actual, err := sut.GetPodCandidatesForNamespace(context.Background(), "orders")
	// THEN
	require.NoError(t, err)
	require.Len(t, actual, 2)
	assert.Contains(t, actual, model.PodCandidate{OwnerName: "statefulset/orders/statefulset-a", Labels: map[string]string{
		"app": "app-a",
	}})
	assert.Contains(t, actual, model.PodCandidate{OwnerName: "statefulset/orders/statefulset-b", Labels: map[string]string{
		"app": "app-b",
	}})
}

func fixStatefulsetA() *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "statefulset-a",
			Namespace: "orders",
		},
		Spec: appsv1.StatefulSetSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "app-a",
					},
				},
			},
		},
	}
}

func fixStatefulsetB() *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "statefulset-b",
			Namespace: "orders",
		},
		Spec: appsv1.StatefulSetSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "app-b",
					},
				},
			},
		},
	}
}
