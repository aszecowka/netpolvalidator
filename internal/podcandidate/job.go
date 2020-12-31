package podcandidate

import (
	"context"
	"fmt"

	v13 "k8s.io/api/batch/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/client-go/kubernetes/typed/batch/v1"

	"github.com/aszecowka/netpolvalidator/internal/model"
)

type JobFetcher struct {
	client v12.JobsGetter
}

func NewJobFetcher(client v12.JobsGetter) *JobFetcher {
	return &JobFetcher{client: client}
}

func (jf *JobFetcher) GetPodCandidatesForNamespace(ctx context.Context, ns string) ([]model.PodCandidate, error) {
	var allJobs []v13.Job
	continueOption := ""
	for {
		list, err := jf.client.Jobs(ns).List(ctx, v1.ListOptions{Continue: continueOption})
		if err != nil {
			return nil, fmt.Errorf("while getting jobs for namespace %s: %w", ns, err)
		}
		allJobs = append(allJobs, list.Items...)
		continueOption = list.Continue
		if continueOption == "" {
			break
		}
	}

	var out []model.PodCandidate
	for _, j := range allJobs {
		out = append(out, jf.convert(j))
	}

	return out, nil
}

func (jf *JobFetcher) convert(job v13.Job) model.PodCandidate {
	// TODO take into account owner ref
	return model.PodCandidate{
		Labels:    job.Spec.Template.Labels,
		OwnerName: getOwnerName(WorkloadJob, job.Namespace, job.Name),
	}
}
