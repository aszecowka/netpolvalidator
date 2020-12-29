package rule

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	v1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/aszecowka/netpolvalidator/internal/model"
)

type labelCorrectness struct{}

func NewLabelCorrectness() *labelCorrectness {
	return &labelCorrectness{}
}

func (lc *labelCorrectness) Validate(state model.ClusterState) ([]model.Violation, error) {
	var resultErr error
	for _, policies := range state.NetworkPolicies {
		for _, np := range policies {
			if err := lc.validateNetworkPolicy(np, state.Namespaces, state.PodCandidates); err != nil {
				resultErr = multierror.Append(resultErr, err)
			}

		}
	}
	return resultErr
}

func (lc *labelCorrectness) validateNetworkPolicy(np netv1.NetworkPolicy, namespaces []v1.Namespace, podCandidates map[string][]model.PodCandidate) ([]model.Violation,error) {
	var resultErr error
	if err := lc.validatePodSelector(np, podCandidates[np.Namespace]); err != nil {
		resultErr = multierror.Append(resultErr, err)
	}
	if err := lc.validateIngress(np, namespaces, podCandidates); err != nil {
		resultErr = multierror.Append(resultErr, err)
	}

	if err := lc.validateEgress(np, namespaces, podCandidates); err != nil {
		resultErr = multierror.Append(resultErr, err)
	}
	return resultErr
}

func (lc *labelCorrectness) validatePodSelector(np netv1.NetworkPolicy, podCandidates []model.PodCandidate) ([]model.Violation,error) {
	selector, err := metav1.LabelSelectorAsSelector(&np.Spec.PodSelector)
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

func (lc *labelCorrectness) validateIngress(np netv1.NetworkPolicy, namespaces []v1.Namespace, podCandidates map[string][]model.PodCandidate) error {
	var resultErr error
	for _, ingresRule := range np.Spec.Ingress {
		for idx, from := range ingresRule.From {
			if from.PodSelector != nil && from.NamespaceSelector != nil {
				filteredNs, err := lc.getNamespacesMatchingSelector(namespaces, *from.NamespaceSelector)
				if err != nil {
					resultErr = multierror.Append(resultErr, fmt.Errorf("while getting namespaces specified in rule [%d] by selector:%w", idx, err))
					continue
				}
				if len(filteredNs) == 0 {
					resultErr = multierror.Append(resultErr, fmt.Errorf("no namespaces that matches namespaces selector in rule [%d]: %v", idx, *from.NamespaceSelector))
					continue
				}
				podsFromNs := lc.getPodsFromNamespaces(filteredNs, podCandidates)
				matching, err := lc.getPodCandidatesMatchingSelector(*from.PodSelector, podsFromNs)
				if err != nil {
					resultErr = multierror.Append(resultErr, fmt.Errorf("while getting pod candidates that matches pod selector in rule [%d]: %w", idx, err))
					continue
				}
				if len(matching) == 0 {
					resultErr = multierror.Append(resultErr, fmt.Errorf(""))
				}
			} else if from.PodSelector != nil {

			} else if from.NamespaceSelector != nil {

			}
		}
	}
	return resultErr
}

func (lc *labelCorrectness) validateEgress(np netv1.NetworkPolicy, namespaces []v1.Namespace, podCandidates map[string][]model.PodCandidate) error {
	return nil
}

func (lc *labelCorrectness) getNamespacesMatchingSelector(in []v1.Namespace, labelSelector metav1.LabelSelector) ([]v1.Namespace, error) {
	selector, err := metav1.LabelSelectorAsSelector(&labelSelector)
	if err != nil {
		return nil, fmt.Errorf("while creating labels.selector: %w", err)
	}

	var out []v1.Namespace
	for _, ns := range in {
		labelsAsASet := labels.Set(ns.Labels)
		if selector.Matches(labelsAsASet) {
			out = append(out, ns)
		}
	}

	return out, nil
}

func (lc *labelCorrectness) getPodCandidatesMatchingSelector(labelSelector metav1.LabelSelector, podCandidates []model.PodCandidate) ([]model.PodCandidate, error) {
	selector, err := metav1.LabelSelectorAsSelector(&labelSelector)
	if err != nil {
		return nil, fmt.Errorf("while creating labels.selector: %w", err)
	}

	var out []model.PodCandidate
	for _, ns := range podCandidates {
		labelsAsASet := labels.Set(ns.Labels)
		if selector.Matches(labelsAsASet) {
			out = append(out, ns)
		}
	}

	return out, nil
}

func (lc *labelCorrectness) getPodsFromNamespaces(namespaces []v1.Namespace, podCandidates map[string][]model.PodCandidate) []model.PodCandidate {
	var out []model.PodCandidate

	for _, ns := range namespaces {
		pods := podCandidates[ns.Name]
		out = append(out, pods...)
	}
	return out

}
