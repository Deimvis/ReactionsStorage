package http_handlers_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/tests/fake"
	setup "github.com/Deimvis/reactionsstorage/tests/setup"
	"github.com/Deimvis/reactionsstorage/tests/simulation/utils"
)

func TestAvailableReactionsGET_Simple(t *testing.T) {
	var req models.AvailableReactionsGETRequest
	req.Query.NamespaceId = fake.NamespaceId

	var resp models.AvailableReactionsGETResponse200
	resp.Reactions = []models.Reaction{
		fake.Reaction, fake.Reaction2, fake.Reaction3,
		fake.FreeReaction1, fake.FreeReaction2, fake.FreeReaction3,
		fake.FreeReaction4, fake.FreeReaction5, fake.FreeReaction6,
	}

	test(t, &req, &resp)
}

func TestAvailableReactionsGET_Complex(t *testing.T) {
	defer setup.SetFakeConfiguration()

	for i, tc := range configurationPOST_testCases_Complex {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			conf := tc.Configuration
			setup.SetConfiguration(conf)
			for _, n := range conf.Namespaces {
				reactionSet := conf.ReactionSets[slices.IndexFunc(conf.ReactionSets, func(rs models.ReactionSet) bool { return rs.Id == n.ReactionSetId })]
				availableReactions := utils.Filter(conf.Reactions, func(r models.Reaction) bool { return !utils.Contains(reactionSet.ReactionIds, r.Id) })

				var req models.AvailableReactionsGETRequest
				req.Query.NamespaceId = fake.NamespaceId

				var resp models.AvailableReactionsGETResponse200
				resp.Reactions = availableReactions

				test(t, &req, &resp)
			}
		})
	}
}

func TestAvailableReactionsGET_404(t *testing.T) {
	defer setup.SetFakeConfiguration()

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
			var req models.AvailableReactionsGETRequest
			req.Query.NamespaceId = tc.namespaceId
			resp := request(t, &req)
			require.Equal(t, 404, resp.Code)
		})
	}
}
