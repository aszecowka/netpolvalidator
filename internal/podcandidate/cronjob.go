package podcandidate

import (
	"context"
	"fmt"

	v1beta12 "k8s.io/api/batch/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/batch/v1beta1"

	"github.com/aszecowka/netpolvalidator/internal/model"
)

type CronjobFetcher struct {
	client v1beta1.CronJobsGetter
}

func NewCronjobFetcher(client v1beta1.CronJobsGetter) *CronjobFetcher {
	return &CronjobFetcher{client: client}
}

func (cf *CronjobFetcher) GetPodCandidatesForNamespace(ctx context.Context, ns string) ([]model.PodCandidate, error) {
	var allCronjobs []v1beta12.CronJob
	continueOption := ""
	for {
		cronjobs, err := cf.client.CronJobs(ns).List(ctx, v1.ListOptions{Continue: continueOption})
		if err != nil {
			return nil, fmt.Errorf("while getting cronjobs for namespace %s: %w", ns, err)
		}
		allCronjobs = append(allCronjobs, cronjobs.Items...)
		continueOption = cronjobs.Continue
		if continueOption == "" {
			break
		}
	}

	var out []model.PodCandidate
	for _, cj := range allCronjobs {
		out = append(out, cf.convert(cj))
	}

	return out, nil
}

func (cf *CronjobFetcher) convert(cronjob v1beta12.CronJob) model.PodCandidate {
	return model.PodCandidate{
		Labels:    cronjob.Spec.JobTemplate.Spec.Template.Labels,
		OwnerName: getOwnerName(WorkloadCronjob, cronjob.Namespace, cronjob.Name),
	}
}
