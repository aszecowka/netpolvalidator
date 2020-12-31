package podcandidate_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/apps/v1"
	v13 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/aszecowka/netpolvalidator/internal/model"
	"github.com/aszecowka/netpolvalidator/internal/podcandidate"
)

func TestPodCandidatesFromDaemonsets(t *testing.T) {
	// GIVEN
	fakeClientset := fake.NewSimpleClientset(fixDsA(), fixDsB())
	sut := podcandidate.NewDaemonsetFetcher(fakeClientset.AppsV1())
	// WHEN
	actual, err := sut.GetPodCandidatesForNamespace(context.Background(), "orders")
	// THEN
	require.NoError(t, err)
	require.Len(t, actual, 2)
	assert.Contains(t, actual, model.PodCandidate{
		OwnerName: "daemonset/orders/ds-a",
		Labels: map[string]string{
			"app": "app-a",
		},
	})
	assert.Contains(t, actual, model.PodCandidate{
		OwnerName: "daemonset/orders/ds-b",
		Labels: map[string]string{
			"app": "app-b",
		},
	})
}

func fixDsA() *v1.DaemonSet {
	return &v1.DaemonSet{
		ObjectMeta: v12.ObjectMeta{
			Name:      "ds-a",
			Namespace: "orders",
		},
		Spec: v1.DaemonSetSpec{
			Template: v13.PodTemplateSpec{
				ObjectMeta: v12.ObjectMeta{
					Labels: map[string]string{
						"app": "app-a",
					},
				},
			},
		},
	}
}

func fixDsB() *v1.DaemonSet {
	return &v1.DaemonSet{
		ObjectMeta: v12.ObjectMeta{
			Name:      "ds-b",
			Namespace: "orders",
		},
		Spec: v1.DaemonSetSpec{
			Template: v13.PodTemplateSpec{
				ObjectMeta: v12.ObjectMeta{
					Labels: map[string]string{
						"app": "app-b",
					},
				},
			},
		},
	}
}
