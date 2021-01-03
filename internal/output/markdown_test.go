package output_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/aszecowka/netpolvalidator/internal/model"
	"github.com/aszecowka/netpolvalidator/internal/output"
)

func TestGenerateMarkdownReport(t *testing.T) {
	sut := output.NewMarkdown()
	actual, err := sut.Generate(context.Background(), model.ClusterState{}, []model.Violation{
		{
			Namespace:         "orders",
			Type:              model.ViolationInvalidLabel,
			NetworkPolicyName: "ingress-all",
			Message:           "something went wrong"},
	})
	require.NoError(t, err)
	require.NotNil(t, actual)

	fmt.Println(actual)
}
