package rule

import (
	"fmt"

	v1 "k8s.io/api/networking/v1"
)

func prettyNetworkPolicy(np v1.NetworkPolicy) string {
	return fmt.Sprintf("%s/%s", np.Namespace, np.Name)
}
