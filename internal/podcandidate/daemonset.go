package podcandidate

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"

	"github.com/aszecowka/netpolvalidator/internal/model"
)

type DaemonsetFetcher struct {
	client v1.DaemonSetsGetter
}

func NewDaemonsetFetcher(client v1.DaemonSetsGetter) *DaemonsetFetcher {
	return &DaemonsetFetcher{client: client}
}

func (df *DaemonsetFetcher) GetPodCandidatesForNamespace(ctx context.Context, ns string) ([]model.PodCandidate, error) {
	var allDs []appsv1.DaemonSet
	continueOption := ""
	for {
		dsList, err := df.client.DaemonSets(ns).List(ctx, metav1.ListOptions{Continue: continueOption})
		if err != nil {
			return nil, fmt.Errorf("while gettting daemonsets from namespace: %s: %w", ns, err)
		}
		allDs = append(allDs, dsList.Items...)
		continueOption = dsList.Continue
		if continueOption == "" {
			break
		}

	}
	var out []model.PodCandidate
	for _, d := range allDs {
		out = append(out, df.convert(d))
	}
	return out, nil
}

func (df *DaemonsetFetcher) convert(daemonset appsv1.DaemonSet) model.PodCandidate {
	return model.PodCandidate{
		Labels:    daemonset.Spec.Template.Labels,
		OwnerName: getOwnerName(WorkloadDaemonset, daemonset.Namespace, daemonset.Name),
	}
}
