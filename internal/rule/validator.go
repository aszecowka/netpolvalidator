package rule

import "github.com/aszecowka/netpolvalidator/internal/model"

type Validator interface {
	Validate(state model.ClusterState) ([]model.Violation, error)
}
