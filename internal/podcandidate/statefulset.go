package podcandidate

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"

	"github.com/aszecowka/netpolvalidator/internal/model"
)

type StatefulsetFetcher struct {
	client v1.StatefulSetsGetter
}

func NewStatefulsetsFetcher(client v1.StatefulSetsGetter) *StatefulsetFetcher {
	return &StatefulsetFetcher{client: client}
}

func (sf *StatefulsetFetcher) GetPodCandidatesForNamespace(ctx context.Context, ns string) ([]model.PodCandidate, error) {
	var allStatefulsets []appsv1.StatefulSet
	continueOption := ""
	for {
		list, err := sf.client.StatefulSets(ns).List(ctx, metav1.ListOptions{Continue: continueOption})
		if err != nil {
			return nil, fmt.Errorf("while gettting statefulsets from namespace: %s: %w", ns, err)
		}
		allStatefulsets = append(allStatefulsets, list.Items...)
		continueOption = list.Continue
		if continueOption == "" {
			break
		}

	}
	var out []model.PodCandidate
	for _, d := range allStatefulsets {
		out = append(out, sf.convert(d))
	}
	return out, nil
}

func (sf *StatefulsetFetcher) convert(ss appsv1.StatefulSet) model.PodCandidate {
	return model.PodCandidate{
		Labels:    ss.Spec.Template.Labels,
		OwnerName: getOwnerName(WorkloadStatefulset, ss.Namespace, ss.Name),
	}
}
