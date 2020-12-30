package rule

import (
	"fmt"

	"github.com/aszecowka/netpolvalidator/internal/model"
)

const (
	msgNoPodsMatchingPodSelector                            = "no pods matching pod selector"
	msgNoNsMatchingLabelsForIngressRulePattern              = "no namespaces matching labels for %s rule [%s]"
	msgNoPodsMatchingLabelsForIngressRulePattern            = "no pods matching labels for %s rule [%s]"
	msgNoPodsInNamespaceMatchingLabelsForIngressRulePattern = "no pods in namespaces matching labels for %s rule: [%s]"
)

func getViolationMessageWithTypeAndPosition(pattern string, ruleType model.RuleType, position string) string {
	return fmt.Sprintf(pattern, ruleType, position)
}
