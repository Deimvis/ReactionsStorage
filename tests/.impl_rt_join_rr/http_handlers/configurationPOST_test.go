package http_handlers_test

import (
	"fmt"
	"net/http"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/src/utils"
	"github.com/Deimvis/reactionsstorage/tests/fake"
	setup "github.com/Deimvis/reactionsstorage/tests/setup"
)

func TestConfigurationPOST_Simple(t *testing.T) {
	defer setup.SetFakeConfiguration()

	testCases := []struct {
		ContentType string
	}{
		{
			ContentType: "application/json",
		},
		{
			ContentType: "application/yaml",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var req models.ConfigurationPOSTRequest
			req.Headers = make(http.Header)
			req.Headers.Set("Content-Type", tc.ContentType)
			var resp models.ConfigurationPOSTResponse200
			resp.Status = "ok"
			test(t, &req, &resp)
		})
	}

}

func TestConfigurationPOST_Complex(t *testing.T) {
	defer setup.SetFakeConfiguration()

	for i, tc := range configurationPOST_testCases_Complex {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var req models.ConfigurationPOSTRequest
			req.Headers = make(http.Header)
			req.Headers.Set("Content-Type", "application/json")
			req.Body = tc.Configuration
			var resp models.ConfigurationPOSTResponse200
			resp.Status = "ok"
			test(t, &req, &resp)
		})
	}
}

func TestConfigurationPOST_ComplexWithDeps(t *testing.T) {
	defer setup.SetFakeConfiguration()

	for i, tc := range configurationPOST_testCases_Complex {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var req models.ConfigurationPOSTRequest
			req.Headers = make(http.Header)
			req.Headers.Set("Content-Type", "application/json")
			req.Body = tc.Configuration
			var resp models.ConfigurationPOSTResponse200
			resp.Status = "ok"
			test(t, &req, &resp)
			checkConfigurationIsApplied(t, tc.Configuration)
		})
	}
}

func TestConfigurationPOST_Invalid(t *testing.T) {
	defer setup.SetFakeConfiguration()

	for i, tc := range configurationPOST_testCases_Invalid {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var req models.ConfigurationPOSTRequest
			req.Headers = make(http.Header)
			req.Headers.Set("Content-Type", "application/json")
			req.Body = tc.Configuration
			resp := request(t, &req)
			require.Equal(t, 422, resp.Code)
		})
	}
}

func checkConfigurationIsApplied(t *testing.T, c models.Configuration) {
	for _, n := range c.Namespaces {
		{
			var req models.NamespaceGETRequest
			req.Query.NamespaceId = n.Id
			var resp models.NamespaceGETResponse200
			resp.Namespace = n
			test(t, &req, &resp)
		}
		{
			reactionSet := c.ReactionSets[slices.IndexFunc(c.ReactionSets, func(rs models.ReactionSet) bool { return rs.Id == n.ReactionSetId })]
			availableReactions := utils.Filter(c.Reactions, func(r models.Reaction) bool { return utils.Contains(reactionSet.ReactionIds, r.Id) })

			var req models.AvailableReactionsGETRequest
			req.Query.NamespaceId = n.Id
			var resp models.AvailableReactionsGETResponse200
			resp.Reactions = availableReactions
			test(t, &req, &resp)
		}
	}
}

type configurationPOST_testCase struct {
	models.Configuration
}

var configurationPOST_testCases_Complex = []configurationPOST_testCase{
	{
		models.Configuration{
			Reactions: []models.Reaction{
				fake.Reaction1,
				fake.Reaction2,
				fake.FreeReaction1,
			},
			ReactionSets: []models.ReactionSet{
				{
					Id:          "reaction_set",
					ReactionIds: []string{fake.Reaction1.Id, fake.Reaction2.Id, fake.FreeReaction1.Id},
				},
			},
			Namespaces: []models.Namespace{
				{
					Id:                         "namespace",
					ReactionSetId:              "reaction_set",
					MaxUniqReactions:           1,
					MutuallyExclusiveReactions: [][]string{{fake.Reaction1.Id, fake.Reaction2.Id}},
				},
			},
		},
	},
}

var configurationPOST_testCases_Invalid = []configurationPOST_testCase{
	{
		models.Configuration{
			ReactionSets: []models.ReactionSet{
				{
					Id:          "reaction_set",
					ReactionIds: []string{"reaction id that does not exist"},
				},
			},
		},
	},
	{
		models.Configuration{
			ReactionSets: []models.ReactionSet{
				{
					Id:          "reaction_set1",
					ReactionIds: nil,
				},
				{
					Id:          "reaction_set2",
					ReactionIds: []string{"reaction set above is odd but ok", "this one has reaction ids that does not exist"},
				},
			},
		},
	},
	{
		models.Configuration{
			Namespaces: []models.Namespace{
				{
					Id:            "namespace",
					ReactionSetId: "reaction set id that does not exist",
				},
			},
		},
	},
	{
		models.Configuration{
			ReactionSets: []models.ReactionSet{
				{
					Id: "reaction_set",
				},
			},
			Namespaces: []models.Namespace{
				{
					Id:               "max_uniq_reactions should be non-negative",
					ReactionSetId:    "reaction_set",
					MaxUniqReactions: -1,
				},
			},
		},
	},
	{
		models.Configuration{
			Reactions: []models.Reaction{
				{
					Id: "reaction1",
				},
				{
					Id: "reaction2",
				},
			},
			ReactionSets: []models.ReactionSet{
				{
					Id:          "reaction_set",
					ReactionIds: []string{"reaction1", "reaction2"},
				},
			},
			Namespaces: []models.Namespace{
				{
					Id:                         "namespace",
					ReactionSetId:              "reaction_set",
					MutuallyExclusiveReactions: [][]string{{"reaction1", "reaction id that does not exist"}},
				},
			},
		},
	},
	{
		models.Configuration{
			Reactions: []models.Reaction{
				{
					Id: "reaction1",
				},
				{
					Id: "reaction2",
				},
			},
			ReactionSets: []models.ReactionSet{
				{
					Id:          "reaction_set",
					ReactionIds: []string{"reaction1"},
				},
			},
			Namespaces: []models.Namespace{
				{
					Id:                         "mutually_exclusive_reactions include reaction that is not in the reaction set",
					ReactionSetId:              "reaction_set",
					MutuallyExclusiveReactions: [][]string{{"reaction1", "reaction2"}},
				},
			},
		},
	},
}
