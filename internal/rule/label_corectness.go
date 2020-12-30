package rule

import (
	"fmt"

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
	var allViolations []model.Violation
	for _, policies := range state.NetworkPolicies {
		for _, np := range policies {
			violations, err := lc.validateNetworkPolicy(np, state.Namespaces, state.PodCandidates)
			if err != nil {
				return nil, err
			}
			allViolations = append(allViolations, violations...)
		}
	}
	return allViolations, nil
}

func (lc *labelCorrectness) validateNetworkPolicy(np netv1.NetworkPolicy, namespaces []v1.Namespace, podCandidates map[string][]model.PodCandidate) ([]model.Violation, error) {
	var allViolations []model.Violation
	violations, err := lc.validatePodSelector(np, podCandidates[np.Namespace])
	if err != nil {
		return nil, err
	}
	allViolations = append(allViolations, violations...)

	violations, err = lc.validateIngress(np, namespaces, podCandidates)
	if err != nil {
		return nil, err
	}
	allViolations = append(allViolations, violations...)

	violations, err = lc.validateEgress(np, namespaces, podCandidates)
	if err != nil {
		return nil, err
	}
	return allViolations, nil
}

func (lc *labelCorrectness) validatePodSelector(np netv1.NetworkPolicy, podCandidates []model.PodCandidate) ([]model.Violation, error) {
	selector, err := metav1.LabelSelectorAsSelector(&np.Spec.PodSelector)
	if err != nil {
		return nil, fmt.Errorf("while creating lables.selector: %w", err)
	}

	found := false
	for _, pc := range podCandidates {
		labelsAsASet := labels.Set(pc.Labels)
		if selector.Matches(labelsAsASet) {
			found = true
			break
		}

	}
	if found {
		return nil, nil
	}
	return []model.Violation{
		model.NewViolation(np, msgNoPodsMatchingPodSelector, model.ViolationInvalidLabel),
	}, nil
}

func (lc *labelCorrectness) validateIngress(np netv1.NetworkPolicy, namespaces []v1.Namespace, podCandidates map[string][]model.PodCandidate) ([]model.Violation, error) {
	var allViolations []model.Violation
	for idxIngress, ingresRule := range np.Spec.Ingress {
		for idxFrom, from := range ingresRule.From {
			position := fmt.Sprintf("%d:%d", idxIngress+1, idxFrom+1)
			if from.PodSelector != nil && from.NamespaceSelector != nil {
				filteredNs, err := lc.getNamespacesMatchingSelector(namespaces, *from.NamespaceSelector)
				if err != nil {
					return nil, fmt.Errorf("while getting namespaces specified in the ingress rule [%s] for [%s/%s]:%w", position, np.Namespace, np.Name, err)
				}
				if len(filteredNs) == 0 {
					allViolations = append(allViolations, model.NewViolation(np, fmt.Sprintf(msgNoNsMatchingLabelsForIngressRulePattern, position), model.ViolationInvalidLabel))
					continue
				}
				podsFromNs := lc.getPodsFromNamespaces(filteredNs, podCandidates)
				matching, err := lc.getPodCandidatesMatchingSelector(*from.PodSelector, podsFromNs)
				if err != nil {
					return nil, fmt.Errorf("while getting pod candidates that matches pod selector in rule [%s] for [%s/%s]: %w", position, np.Namespace, np.Name, err)

				}
				if len(matching) == 0 {
					allViolations = append(allViolations, model.NewViolation(np, fmt.Sprintf(msgNoPodsMatchingLabelsForIngressRulePattern, position), model.ViolationInvalidLabel))
					continue
				}
			} else if from.PodSelector != nil {
				podsInTheSameNs, err := lc.getPodCandidatesMatchingSelector(*from.PodSelector, podCandidates[np.Namespace])
				if err != nil {
					return nil, fmt.Errorf("while getting pod candidates that matches pod selector in rule [%s] for [%s/%s]: %w", position, np.Namespace, np.Name, err)
				}
				if len(podsInTheSameNs) == 0 {
					allViolations = append(allViolations, model.NewViolation(np, fmt.Sprintf(msgNoPodsMatchingLabelsForIngressRulePattern, position), model.ViolationInvalidLabel))
					continue
				}
			} else if from.NamespaceSelector != nil {
				filteredNs, err := lc.getNamespacesMatchingSelector(namespaces, *from.NamespaceSelector)
				if err != nil {
					return nil, fmt.Errorf("while getting namespaces specified in the ingress rule [%s] for [%s/%s]:%w", position, np.Namespace, np.Name, err)
				}
				if len(filteredNs) == 0 {
					allViolations = append(allViolations, model.NewViolation(np, fmt.Sprintf(msgNoNsMatchingLabelsForIngressRulePattern, position), model.ViolationInvalidLabel))
					continue
				}

				podsInFilteredNS := 0
				for _, ns := range filteredNs {
					podsInFilteredNS += len(podCandidates[ns.Name])
				}
				if podsInFilteredNS == 0 {
					allViolations = append(allViolations, model.NewViolation(np, fmt.Sprintf(msgNoPodsInNamespaceMatchingLabelsForIngressRulePattern, position), model.ViolationInvalidLabel))
					continue
				}
			}
		}
	}
	return allViolations, nil
}

func (lc *labelCorrectness) validateEgress(np netv1.NetworkPolicy, namespaces []v1.Namespace, podCandidates map[string][]model.PodCandidate) ([]model.Violation, error) {
	return nil, nil
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
