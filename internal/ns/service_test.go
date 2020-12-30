package ns_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/aszecowka/netpolvalidator/internal/ns"
)

func TestGetAllNamespaces(t *testing.T) {
	// GIVEN
	nsA := fixNsA()
	nsB := fixNsB()
	fakeClientset := fake.NewSimpleClientset(&nsA, &nsB)
	sut := ns.New(fakeClientset.CoreV1().Namespaces())
	// WHEN
	actual, err := sut.GetAllNamespaces(context.Background())
	// THEN
	require.NoError(t, err)
	require.Len(t, actual, 2)
	assert.Contains(t, actual, nsA)
	assert.Contains(t, actual, nsB)

}
func fixNsA() v1.Namespace {
	return v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "a",
		},
	}
}

func fixNsB() v1.Namespace {
	return v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "b",
		},
	}
}
