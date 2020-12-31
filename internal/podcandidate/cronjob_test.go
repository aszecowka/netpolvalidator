package podcandidate_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v12 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	v13 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/aszecowka/netpolvalidator/internal/model"
	"github.com/aszecowka/netpolvalidator/internal/podcandidate"
)

func TestPodCandidatesFromCronjobs(t *testing.T) {
	// GIVEN
	crA := fixCronJobA()
	crB := fixCronJobB()
	fakeClientset := fake.NewSimpleClientset(&crA, &crB)
	sut := podcandidate.NewCronjobFetcher(fakeClientset.BatchV1beta1())
	// WHEN
	actual, err := sut.GetPodCandidatesForNamespace(context.Background(), "orders")
	// THEN
	require.NoError(t, err)
	require.Len(t, actual, 2)
	assert.Contains(t, actual, model.PodCandidate{
		OwnerName: "cronjob/orders/cron-job-a",
		Labels: map[string]string{
			"app": "app-a",
		},
	})
	assert.Contains(t, actual, model.PodCandidate{
		OwnerName: "cronjob/orders/cron-job-b",
		Labels: map[string]string{
			"app": "app-b",
		},
	})
}

func fixCronJobA() v1beta1.CronJob {
	return v1beta1.CronJob{
		ObjectMeta: v1.ObjectMeta{
			Name:      "cron-job-a",
			Namespace: "orders",
		},
		Spec: v1beta1.CronJobSpec{
			JobTemplate: v1beta1.JobTemplateSpec{
				Spec: v12.JobSpec{
					Template: v13.PodTemplateSpec{
						ObjectMeta: v1.ObjectMeta{
							Labels: map[string]string{
								"app": "app-a",
							},
						},
					},
				},
			},
		},
	}
}

func fixCronJobB() v1beta1.CronJob {
	return v1beta1.CronJob{
		ObjectMeta: v1.ObjectMeta{
			Name:      "cron-job-b",
			Namespace: "orders",
		},
		Spec: v1beta1.CronJobSpec{
			JobTemplate: v1beta1.JobTemplateSpec{
				Spec: v12.JobSpec{
					Template: v13.PodTemplateSpec{
						ObjectMeta: v1.ObjectMeta{
							Labels: map[string]string{
								"app": "app-b",
							},
						},
					},
				},
			},
		},
	}
}
