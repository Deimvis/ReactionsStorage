package http_handlers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/tests/fake"
	setup "github.com/Deimvis/reactionsstorage/tests/setup"
)

func TestReactionsDELETE_Simple(t *testing.T) {
	defer setup.ClearUserReactions()

	var req models.ReactionsDELETERequest
	req.Body.NamespaceId = fake.NamespaceId
	req.Body.EntityId = fake.EntityId
	req.Body.ReactionId = fake.ReactionId
	req.Body.UserId = fake.UserId

	var resp models.ReactionsDELETEResponse200
	resp.Status = "ok"

	test(t, &req, &resp)
}

func TestReactionsDELETE_Complex(t *testing.T) {
	defer setup.ClearUserReactions()

	for _, tc := range reactionsDELETE_testCase_Complex {
		t.Run(tc.name, func(t *testing.T) {
			setup.SetUserReactions(tc.initialURs)
			for _, r := range tc.requests {
				var req models.ReactionsDELETERequest
				req.Body = r.userReaction
				resp := request(t, &req)
				require.Equal(t, r.expectedStatusCode, resp.Code)
			}
			require.ElementsMatch(t, tc.expectedURs, setup.GetAllUserReactions())
		})
	}
}

type reactionsDELETE_testCase struct {
	name        string
	initialURs  []models.UserReaction
	requests    []reactionsDELETE_request
	expectedURs []models.UserReaction
}

type reactionsDELETE_request struct {
	userReaction       models.UserReaction
	expectedStatusCode int
}

var reactionsDELETE_testCase_Complex = []reactionsDELETE_testCase{
	{
		name: "remove 1 reaction",
		initialURs: []models.UserReaction{
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId,
				ReactionId:  fake.ReactionId,
			},
		},
		requests: []reactionsDELETE_request{
			{
				userReaction: models.UserReaction{
					NamespaceId: fake.NamespaceId,
					EntityId:    fake.EntityId,
					UserId:      fake.UserId,
					ReactionId:  fake.ReactionId,
				},
				expectedStatusCode: 200,
			},
		},
		expectedURs: []models.UserReaction{},
	},
	{
		name:       "remove 1 reaction that does not exist",
		initialURs: nil,
		requests: []reactionsDELETE_request{
			{
				userReaction: models.UserReaction{
					NamespaceId: fake.NamespaceId,
					EntityId:    fake.EntityId,
					UserId:      fake.UserId,
					ReactionId:  fake.ReactionId,
				},
				expectedStatusCode: 200,
			},
		},
		expectedURs: []models.UserReaction{},
	},
	{
		name: "remove 1 reaction out of 2 same reactions",
		initialURs: []models.UserReaction{
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
		requests: []reactionsDELETE_request{
			{
				userReaction: models.UserReaction{
					NamespaceId: fake.NamespaceId,
					EntityId:    fake.EntityId,
					UserId:      fake.UserId1,
					ReactionId:  fake.ReactionId,
				},
				expectedStatusCode: 200,
			},
		},
		expectedURs: []models.UserReaction{
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId2,
				ReactionId:  fake.ReactionId,
			},
		},
	},
	{
		name: "remove 1 reaction out of 2 different reactions",
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
		},
		requests: []reactionsDELETE_request{
			{
				userReaction: models.UserReaction{
					NamespaceId: fake.NamespaceId,
					EntityId:    fake.EntityId,
					UserId:      fake.UserId,
					ReactionId:  fake.FreeReactionId1,
				},
				expectedStatusCode: 200,
			},
		},
		expectedURs: []models.UserReaction{
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId,
				ReactionId:  fake.FreeReactionId2,
			},
		},
	},
}
