package http_handlers_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/tests/fake"
	setup "github.com/Deimvis/reactionsstorage/tests/setup"
)

func TestNamespaceGET_Simple(t *testing.T) {
	setup.SetFakeConfiguration()

	var req models.NamespaceGETRequest
	req.Query.NamespaceId = fake.NamespaceId

	var resp models.NamespaceGETResponse200
	resp.Namespace = fake.Namespace

	test(t, &req, &resp)
}

func TestNamespaceGET_Complex(t *testing.T) {
	defer setup.SetFakeConfiguration()

	for i, tc := range namespaceGET_testCases_Complex {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			setup.SetConfiguration(models.Configuration{Namespaces: tc.Namespaces})
			for _, n := range tc.Namespaces {
				var req models.NamespaceGETRequest
				req.Query.NamespaceId = n.Id
				var resp models.NamespaceGETResponse200
				resp.Namespace = n
				test(t, &req, &resp)
			}
		})
	}
}

func TestNamespaceGET_404(t *testing.T) {
	setup.SetFakeConfiguration()

	testCases := []struct {
		namespaceId string
	}{
		{
			"namespace that does not exist",
		},
		{
			"namespace that does not exist 2",
		},
		{
			"namespace that does not exist 3",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var req models.NamespaceGETRequest
			req.Query.NamespaceId = tc.namespaceId
			resp := request(t, &req)
			require.Equal(t, 404, resp.Code)
		})
	}
}

type namespaceGET_testCase struct {
	Namespaces []models.Namespace
}

var namespaceGET_testCases_Complex = []namespaceGET_testCase{
	{
		Namespaces: []models.Namespace{
			{
				Id:                         "namespace",
				MutuallyExclusiveReactions: [][]string{},
			},
		},
	},
	{
		Namespaces: []models.Namespace{
			{
				Id:                         "namespace1",
				MutuallyExclusiveReactions: [][]string{},
			},
			{
				Id:                         "namespace2",
				MutuallyExclusiveReactions: [][]string{},
			},
		},
	},
}
