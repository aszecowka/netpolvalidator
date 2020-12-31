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

func TestPodCandidateFromDeployments(t *testing.T) {
	// GIVEN
	deployA := fixDeployA()
	deployB := fixDeployB()
	fakeClientset := fake.NewSimpleClientset(&deployA, &deployB)
	sut := podcandidate.NewDeploymentsFetcher(fakeClientset.AppsV1())
	// WHEN
	actual, err := sut.GetPodCandidatesForNamespace(context.Background(), "orders")
	// THEN
	require.NoError(t, err)
	require.Len(t, actual, 2)
	assert.Contains(t, actual, model.PodCandidate{OwnerName: "deployment/orders/deploy-a", Labels: map[string]string{
		"app": "app-a",
	}})
	assert.Contains(t, actual, model.PodCandidate{OwnerName: "deployment/orders/deploy-b", Labels: map[string]string{
		"app": "app-b",
	}})
}

func fixDeployA() appsv1.Deployment {
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deploy-a",
			Namespace: "orders",
		},
		Spec: appsv1.DeploymentSpec{
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

func fixDeployB() appsv1.Deployment {
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deploy-b",
			Namespace: "orders",
		},
		Spec: appsv1.DeploymentSpec{
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
