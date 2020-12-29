package podcandidate

import (
	"context"
	"fmt"

	v13 "k8s.io/api/apps/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"

	"github.com/aszecowka/netpolvalidator/internal/model"
)

type DeployService struct {
	client v1.DeploymentsGetter
}

func NewDeploymentService(client v1.DeploymentsGetter) *DeployService {
	return &DeployService{client: client}
}

func (ds *DeployService) GetPodCandidatesForNamespace(ctx context.Context, ns string) ([]model.PodCandidate, error) {
	deployments, err := ds.client.Deployments(ns).List(ctx, v12.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("while gettting deployments from namespace: %s: %w", ns, err)
	}
	var out []model.PodCandidate
	for _, d := range deployments.Items {
		out = append(out, ds.convert(d))
	}
	return out, nil
}

func (ds *DeployService) convert(deploy v13.Deployment) model.PodCandidate {
	return model.PodCandidate{
		Labels:    deploy.Spec.Template.Labels,
		OwnerName: fmt.Sprintf("deploy/%s%s", deploy.Namespace, deploy.Name),
	}
}