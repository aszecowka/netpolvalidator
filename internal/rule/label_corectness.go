package rule

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	v1 "k8s.io/api/networking/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/aszecowka/netpolvalidator/internal/model"
)

type labelCorrectness struct{}

func NewLabelCorrectness() *labelCorrectness {
	return &labelCorrectness{}
}

func (lc *labelCorrectness) Validate(state model.ClusterState) error {
	var resultErr error
	for ns, policies := range state.NetworkPolicies {
		for _, np := range policies {
			if err := lc.validateNetworkPolicy(ns, np, state.PodCandidates); err != nil {
				resultErr = multierror.Append(resultErr, err)
			}

		}
	}
	return resultErr
}

func (lc *labelCorrectness) validateNetworkPolicy(ns string, np v1.NetworkPolicy, podCandidates map[string][]model.PodCandidate) error {
	if err := lc.validatePodSelector(np, podCandidates[ns]); err != nil {
		return err
	}
	return nil
}

func (lc *labelCorrectness) validatePodSelector(np v1.NetworkPolicy, podCandidates []model.PodCandidate) error {
	selector, err := v12.LabelSelectorAsSelector(&np.Spec.PodSelector)
	if err != nil {
		return fmt.Errorf("while creating lables.selector: %w", err)
	}

	found := false
	for _, pc := range podCandidates {
		labelsAsASet := labels.Set(pc.Labels)
		if selector.Matches(labelsAsASet) {
			found = true
			break
		}

	}
	if !found {
		return fmt.Errorf("there is no matching pods for network policy: [%s/%s]", np.Namespace, np.Name)
	}
	return nil
}
