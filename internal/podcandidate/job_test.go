package podcandidate_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v12 "k8s.io/api/batch/v1"
	v13 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/aszecowka/netpolvalidator/internal/model"
	"github.com/aszecowka/netpolvalidator/internal/podcandidate"
)

func TestPodCandidatesFromJobs(t *testing.T) {
	// GIVEN
	fakeClientset := fake.NewSimpleClientset(fixJobA(), fixJobB())
	sut := podcandidate.NewJobFetcher(fakeClientset.BatchV1())
	// WHEN
	actual, err := sut.GetPodCandidatesForNamespace(context.Background(), "orders")
	// THEN
	require.NoError(t, err)
	require.Len(t, actual, 2)
	assert.Contains(t, actual, model.PodCandidate{
		OwnerName: "job/orders/job-a",
		Labels: map[string]string{
			"app": "app-a",
		},
	})
	assert.Contains(t, actual, model.PodCandidate{
		OwnerName: "job/orders/job-b",
		Labels: map[string]string{
			"app": "app-b",
		},
	})
}

func fixJobA() *v12.Job {
	return &v12.Job{
		ObjectMeta: v1.ObjectMeta{
			Name:      "job-a",
			Namespace: "orders",
		},
		Spec: v12.JobSpec{
			Template: v13.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						"app": "app-a",
					},
				},
			},
		},
	}
}

func fixJobB() *v12.Job {
	return &v12.Job{
		ObjectMeta: v1.ObjectMeta{
			Name:      "job-b",
			Namespace: "orders",
		},
		Spec: v12.JobSpec{
			Template: v13.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						"app": "app-b",
					},
				},
			},
		},
	}
}
