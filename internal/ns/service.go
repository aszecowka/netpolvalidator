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

func (s *service) GetNamespacesBySelector(ctx context.Context, selector metaV1.LabelSelector) ([]v1.Namespace, error) {
	// TODO: continue field is ignored
	list, err := s.client.List(ctx, metaV1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		return nil, fmt.Errorf("while listing namespaces by selector: %w", err)
	}

	return list.Items, nil
}

func (s *service) GetAllNamespaces(ctx context.Context) ([]v1.Namespace, error) {
	// TODO: continue field is ignored
	list, err := s.client.List(ctx, metaV1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("while listing all namespaces: %w", err)
	}

	return list.Items, nil
}
