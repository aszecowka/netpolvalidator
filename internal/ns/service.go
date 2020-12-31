package ns

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func New(client typedCoreV1.NamespaceInterface) *service {
	return &service{client: client}
}

type service struct {
	client typedCoreV1.NamespaceInterface
}

func (s *service) GetAllNamespaces(ctx context.Context) ([]v1.Namespace, error) {
	var response []v1.Namespace
	continueOption := ""

	for {
		list, err := s.client.List(ctx, metaV1.ListOptions{Continue: continueOption})
		if err != nil {
			return nil, fmt.Errorf("while listing all namespaces: %w", err)
		}
		response = append(response, list.Items...)
		continueOption = list.Continue
		if continueOption == "" {
			break
		}
	}

	return response, nil
}
