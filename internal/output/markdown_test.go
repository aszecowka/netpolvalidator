package output_test

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aszecowka/netpolvalidator/internal/model"
	"github.com/aszecowka/netpolvalidator/internal/output"
)

func TestGenerateMarkdownReport(t *testing.T) {
	sut := output.NewMarkdown()

	t.Run("no violations", func(t *testing.T) {
		// WHEN
		actual, err := sut.Generate(context.Background(), model.ClusterState{}, []model.Violation{})
		// THEN
		require.NoError(t, err)
		require.NotNil(t, actual)
		actualBytes, err := ioutil.ReadAll(actual)
		require.NoError(t, err)
		expected := getGoldenFileContent(t, "testdata/no_violations.md")
		assert.Equal(t, expected, string(actualBytes))

	})

	t.Run("single violation", func(t *testing.T) {
		// WHEN
		actual, err := sut.Generate(context.Background(), model.ClusterState{}, []model.Violation{
			{
				Namespace:         "orders",
				Type:              model.ViolationInvalidLabel,
				NetworkPolicyName: "ingress-all",
				Message:           "something went wrong",
			},
		})
		// THEN
		require.NoError(t, err)
		require.NotNil(t, actual)
		actualBytes, err := ioutil.ReadAll(actual)
		require.NoError(t, err)
		expected := getGoldenFileContent(t, "testdata/single_violation.md")
		assert.Equal(t, expected, string(actualBytes))
	})

	t.Run("many violations", func(t *testing.T) {
		// WHEN
		actual, err := sut.Generate(context.Background(), model.ClusterState{}, []model.Violation{
			{
				Namespace:         "orders",
				Type:              model.ViolationInvalidLabel,
				NetworkPolicyName: "ingress-all",
				Message:           "something went wrong",
			},
			{
				Namespace:         "users",
				Type:              model.ViolationInvalidLabel,
				NetworkPolicyName: "egress-all",
				Message:           "big mistake",
			},
		})
		// THEN
		require.NoError(t, err)
		require.NotNil(t, actual)
		actualBytes, err := ioutil.ReadAll(actual)
		require.NoError(t, err)
		expected := getGoldenFileContent(t, "testdata/many_violations.md")
		assert.Equal(t, expected, string(actualBytes))
	})

}

func getGoldenFileContent(t *testing.T, path string) string {
	t.Helper()
	expectedFile, err := os.Open(path)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, expectedFile.Close())
	}()
	expectedBytes, err := ioutil.ReadAll(expectedFile)
	require.NoError(t, err)
	return string(expectedBytes)
}
