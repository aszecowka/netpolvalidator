package podcandidate

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"

	"github.com/aszecowka/netpolvalidator/internal/model"
)

type DeploymentsFetcher struct {
	client v1.DeploymentsGetter
}

func NewDeploymentsFetcher(client v1.DeploymentsGetter) *DeploymentsFetcher {
	return &DeploymentsFetcher{client: client}
}

func (df *DeploymentsFetcher) GetPodCandidatesForNamespace(ctx context.Context, ns string) ([]model.PodCandidate, error) {
	var allDeployments []appsv1.Deployment
	continueOption := ""
	for {
		deployments, err := df.client.Deployments(ns).List(ctx, metav1.ListOptions{Continue: continueOption })
		if err != nil {
			return nil, fmt.Errorf("while gettting deployments from namespace: %s: %w", ns, err)
		}
		allDeployments = append(allDeployments, deployments.Items...)
		continueOption = deployments.Continue
		if continueOption == "" {
			break
		}

	}
	var out []model.PodCandidate
	for _, d := range allDeployments {
		out = append(out, df.convert(d))
	}
	return out, nil
}

func (df *DeploymentsFetcher) convert(deploy appsv1.Deployment) model.PodCandidate {
	return model.PodCandidate{
		Labels:    deploy.Spec.Template.Labels,
		OwnerName: getOwnerName(WorkloadDeployment, deploy.Namespace, deploy.Name),
	}
}
