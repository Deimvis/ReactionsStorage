package http_handlers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/tests/fake"
	setup "github.com/Deimvis/reactionsstorage/tests/setup"
)

func TestReactionsPOST_Simple(t *testing.T) {
	defer setup.ClearUserReactions()

	var req models.ReactionsPOSTRequest
	req.Body.NamespaceId = fake.NamespaceId
	req.Body.EntityId = fake.EntityId
	req.Body.ReactionId = fake.ReactionId
	req.Body.UserId = fake.UserId

	var resp models.ReactionsPOSTResponse200
	resp.Status = "ok"

	test(t, &req, &resp)
}

func TestReactionsPOST_Complex(t *testing.T) {
	defer setup.ClearUserReactions()

	for _, tc := range reactionsPOST_testCase_Complex {
		t.Run(tc.name, func(t *testing.T) {
			setup.SetUserReactions(tc.initialURs)
			for _, r := range tc.requests {
				var req models.ReactionsPOSTRequest
				req.Body = r.userReaction
				req.Query.Force = &r.force
				resp := request(t, &req)
				require.Equal(t, r.expectedStatusCode, resp.Code)
			}
			require.ElementsMatch(t, tc.expectedURs, setup.GetUserReactions())
		})
	}
}

func TestReactionsPOST_403(t *testing.T) {
	defer setup.ClearUserReactions()

	for _, tc := range reactionsPOST_testCase_403 {
		t.Run(tc.name, func(t *testing.T) {
			setup.SetUserReactions(tc.initialURs)
			for _, r := range tc.requests {
				var req models.ReactionsPOSTRequest
				req.Body = r.userReaction
				req.Query.Force = &r.force
				resp := request(t, &req)
				require.Equal(t, 403, resp.Code)
			}
			require.ElementsMatch(t, tc.expectedURs, setup.GetUserReactions())
		})
	}
}

type reactionsPOST_testCase struct {
	name        string
	initialURs  []models.UserReaction
	requests    []reactionsPOST_request
	expectedURs []models.UserReaction
}

type reactionsPOST_request struct {
	userReaction       models.UserReaction
	force              bool
	expectedStatusCode int
}

var reactionsPOST_testCase_Complex = []reactionsPOST_testCase{
	{
		name:       "add 1 reaction",
		initialURs: nil,
		requests: []reactionsPOST_request{
			{
				userReaction: models.UserReaction{
					NamespaceId: fake.NamespaceId,
					EntityId:    fake.EntityId,
					UserId:      fake.UserId,
					ReactionId:  fake.ReactionId,
				},
				force:              true,
				expectedStatusCode: 200,
			},
		},
		expectedURs: []models.UserReaction{
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId,
				ReactionId:  fake.ReactionId,
			},
		},
	},
	{
		name:       "add 2 reactions to different entities",
		initialURs: nil,
		requests: []reactionsPOST_request{
			{
				userReaction: models.UserReaction{
					NamespaceId: fake.NamespaceId,
					EntityId:    fake.EntityId1,
					UserId:      fake.UserId,
					ReactionId:  fake.ReactionId,
				},
				force:              true,
				expectedStatusCode: 200,
			},
			{
				userReaction: models.UserReaction{
					NamespaceId: fake.NamespaceId,
					EntityId:    fake.EntityId2,
					UserId:      fake.UserId,
					ReactionId:  fake.ReactionId,
				},
				force:              true,
				expectedStatusCode: 200,
			},
		},
		expectedURs: []models.UserReaction{
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId1,
				UserId:      fake.UserId,
				ReactionId:  fake.ReactionId,
			},
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId2,
				UserId:      fake.UserId,
				ReactionId:  fake.ReactionId,
			},
		},
	},
	{
		name:       "add 2 reactions from different users",
		initialURs: nil,
		requests: []reactionsPOST_request{
			{
				userReaction: models.UserReaction{
					NamespaceId: fake.NamespaceId,
					EntityId:    fake.EntityId,
					UserId:      fake.UserId1,
					ReactionId:  fake.ReactionId,
				},
				force:              true,
				expectedStatusCode: 200,
			},
			{
				userReaction: models.UserReaction{
					NamespaceId: fake.NamespaceId,
					EntityId:    fake.EntityId,
					UserId:      fake.UserId2,
					ReactionId:  fake.ReactionId,
				},
				force:              true,
				expectedStatusCode: 200,
			},
		},
		expectedURs: []models.UserReaction{
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId1,
				ReactionId:  fake.ReactionId,
			},
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId2,
				ReactionId:  fake.ReactionId,
			},
		},
	},
	{
		name:       "add 2 different reactions",
		initialURs: nil,
		requests: []reactionsPOST_request{
			{
				userReaction: models.UserReaction{
					NamespaceId: fake.NamespaceId,
					EntityId:    fake.EntityId,
					UserId:      fake.UserId,
					ReactionId:  fake.FreeReactionId1,
				},
				force:              true,
				expectedStatusCode: 200,
			},
			{
				userReaction: models.UserReaction{
					NamespaceId: fake.NamespaceId,
					EntityId:    fake.EntityId,
					UserId:      fake.UserId,
					ReactionId:  fake.FreeReactionId2,
				},
				force:              true,
				expectedStatusCode: 200,
			},
		},
		expectedURs: []models.UserReaction{
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId,
				ReactionId:  fake.FreeReactionId1,
			},
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId,
				ReactionId:  fake.FreeReactionId2,
			},
		},
	},
	{
		name: "add 1 reaction over existing from another user",
		initialURs: []models.UserReaction{
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId1,
				ReactionId:  fake.ReactionId,
			},
		},
		requests: []reactionsPOST_request{
			{
				userReaction: models.UserReaction{
					NamespaceId: fake.NamespaceId,
					EntityId:    fake.EntityId,
					UserId:      fake.UserId2,
					ReactionId:  fake.ReactionId,
				},
				force:              true,
				expectedStatusCode: 200,
			},
		},
		expectedURs: []models.UserReaction{
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId1,
				ReactionId:  fake.ReactionId,
			},
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId2,
				ReactionId:  fake.ReactionId,
			},
		},
	},
	{
		name: "replace conflicting reaction on force=true",
		initialURs: []models.UserReaction{
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId,
				ReactionId:  fake.ReactionId1,
			},
		},
		requests: []reactionsPOST_request{
			{
				userReaction: models.UserReaction{
					NamespaceId: fake.NamespaceId,
					EntityId:    fake.EntityId,
					UserId:      fake.UserId,
					ReactionId:  fake.ReactionId2,
				},
				force:              true,
				expectedStatusCode: 200,
			},
		},
		expectedURs: []models.UserReaction{
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId,
				ReactionId:  fake.ReactionId2,
			},
		},
	},
}

var reactionsPOST_testCase_403 = []reactionsPOST_testCase{
	{
		name: "can't replace conflicting reaction on force=false",
		initialURs: []models.UserReaction{
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId,
				ReactionId:  fake.ReactionId1,
			},
		},
		requests: []reactionsPOST_request{
			{
				userReaction: models.UserReaction{
					NamespaceId: fake.NamespaceId,
					EntityId:    fake.EntityId,
					UserId:      fake.UserId,
					ReactionId:  fake.ReactionId2,
				},
				force:              false,
				expectedStatusCode: 403,
			},
		},
		expectedURs: []models.UserReaction{
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId,
				ReactionId:  fake.ReactionId1,
			},
		},
	},
	{
		name: "can't add reaction, because of max uniq reactions limit (5)",
		initialURs: []models.UserReaction{
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId,
				ReactionId:  fake.FreeReactionId1,
			},
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId,
				ReactionId:  fake.FreeReactionId2,
			},
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId,
				ReactionId:  fake.FreeReactionId3,
			},
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId,
				ReactionId:  fake.FreeReactionId4,
			},
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId,
				ReactionId:  fake.FreeReactionId5,
			},
		},
		requests: []reactionsPOST_request{
			{
				userReaction: models.UserReaction{
					NamespaceId: fake.NamespaceId,
					EntityId:    fake.EntityId,
					UserId:      fake.UserId,
					ReactionId:  fake.FreeReactionId6,
				},
				force:              true,
				expectedStatusCode: 403,
			},
		},
		expectedURs: []models.UserReaction{
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId,
				ReactionId:  fake.FreeReactionId1,
			},
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId,
				ReactionId:  fake.FreeReactionId2,
			},
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId,
				ReactionId:  fake.FreeReactionId3,
			},
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId,
				ReactionId:  fake.FreeReactionId4,
			},
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId,
				ReactionId:  fake.FreeReactionId5,
			},
		},
	},
}
