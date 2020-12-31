package podcandidate

import (
	"context"
	"fmt"

	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/aszecowka/netpolvalidator/internal/model"
)

type PodsFetcher struct {
	client v1.PodsGetter
}

func NewPodsFetcher(client v1.PodsGetter) *PodsFetcher {
	return &PodsFetcher{client: client}
}

func (pf *PodsFetcher) GetPodCandidatesForNamespace(ctx context.Context, ns string) ([]model.PodCandidate, error) {
	var allPods []v12.Pod
	continueOption := ""
	for {
		list, err := pf.client.Pods(ns).List(ctx, metav1.ListOptions{Continue: continueOption})
		if err != nil {
			return nil, fmt.Errorf("while gettting pods from namespace: %s: %w", ns, err)
		}
		allPods = append(allPods, list.Items...)
		continueOption = list.Continue
		if continueOption == "" {
			break
		}

	}
	var out []model.PodCandidate
	for _, p := range allPods {
		out = append(out, pf.convert(p))
	}
	return out, nil
}

func (pf *PodsFetcher) convert(pod v12.Pod) model.PodCandidate {
	return model.PodCandidate{
		Labels:    pod.Labels,
		OwnerName: getOwnerName(WorkloadPod, pod.Namespace, pod.Name),
	}
}
