package http_handlers_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/tests/fake"
	"github.com/Deimvis/reactionsstorage/tests/setup"
)

func TestReactionsGET_Simple(t *testing.T) {
	w := requestRaw(t, "GET", "/reactions?namespace_id=namespace&entity_id=entity&user_id=user", nil)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"entity_id":"entity","reactions_count":[],"user_reactions":{"user_id":"user","reactions":[]}}`, w.Body.String())
}

func TestReactionsGET_Complex(t *testing.T) {
	defer setup.ClearUserReactions()

	for _, tc := range reactionsGET_testCase_Complex {
		t.Run(tc.name, func(t *testing.T) {
			setup.SetUserReactions(tc.userReactions)
			var req models.ReactionsGETRequest
			for _, check := range tc.checks {
				req.Query.NamespaceId = check.namespaceId
				req.Query.EntityId = check.entityId
				req.Query.UserId = check.userId

				resp := request(t, &req)
				require.Equal(t, check.resp.Code(), resp.Code)
				if resp.Code != 200 {
					requireResponse(t, check.resp, resp)
				} else {
					// avoid JSONEq since it counts different ordering
					// of json array as an error
					exp := check.resp.(*models.ReactionsGETResponse200)
					act := models.ReactionsGETResponse200{}
					err := json.NewDecoder(resp.Body).Decode(&act)
					require.NoError(t, err)
					require.Equal(t, exp.EntityId, act.EntityId)
					require.ElementsMatch(t, exp.ReactionsCount, act.ReactionsCount)
					require.Equal(t, exp.UserReactions.UserId, act.UserReactions.UserId)
					require.ElementsMatch(t, exp.UserReactions.Reactions, act.UserReactions.Reactions)
				}
			}
		})
	}
}

type reactionsGET_testCase struct {
	name          string
	userReactions []models.UserReaction
	checks        []reactionsGET_check
}

type reactionsGET_check struct {
	namespaceId string
	entityId    string
	userId      string
	resp        models.Response
}

var reactionsGET_testCase_Complex = []reactionsGET_testCase{
	{
		name:          "no reactions",
		userReactions: nil,
		checks: []reactionsGET_check{
			{
				fake.NamespaceId,
				fake.EntityId,
				fake.UserId,
				&models.ReactionsGETResponse200{
					EntityId:       fake.EntityId,
					ReactionsCount: []models.ReactionCount{},
					UserReactions: models.UserReactionsWithinEntity{
						UserId:    fake.UserId,
						Reactions: []string{},
					},
				},
			},
		},
	},
	{
		name: "1 reaction",
		userReactions: []models.UserReaction{
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId,
				ReactionId:  fake.ReactionId,
			},
		},
		checks: []reactionsGET_check{
			{
				fake.NamespaceId,
				fake.EntityId,
				fake.UserId,
				&models.ReactionsGETResponse200{
					EntityId: fake.EntityId,
					ReactionsCount: []models.ReactionCount{
						{
							ReactionId: fake.ReactionId,
							Count:      1,
						},
					},
					UserReactions: models.UserReactionsWithinEntity{
						UserId:    fake.UserId,
						Reactions: []string{fake.ReactionId},
					},
				},
			},
		},
	},
	{
		name: "1 reaction from another user",
		userReactions: []models.UserReaction{
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId1,
				ReactionId:  fake.ReactionId,
			},
		},
		checks: []reactionsGET_check{
			{
				fake.NamespaceId,
				fake.EntityId,
				fake.UserId2,
				&models.ReactionsGETResponse200{
					EntityId: fake.EntityId,
					ReactionsCount: []models.ReactionCount{
						{
							ReactionId: fake.ReactionId,
							Count:      1,
						},
					},
					UserReactions: models.UserReactionsWithinEntity{
						UserId:    fake.UserId2,
						Reactions: []string{},
					},
				},
			},
		},
	},
	{
		name: "1 reaction on another post",
		userReactions: []models.UserReaction{
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId1,
				UserId:      fake.UserId,
				ReactionId:  fake.ReactionId,
			},
		},
		checks: []reactionsGET_check{
			{
				fake.NamespaceId,
				fake.EntityId2,
				fake.UserId,
				&models.ReactionsGETResponse200{
					EntityId:       fake.EntityId2,
					ReactionsCount: []models.ReactionCount{},
					UserReactions: models.UserReactionsWithinEntity{
						UserId:    fake.UserId,
						Reactions: []string{},
					},
				},
			},
		},
	},
	{
		name: "same reactions from other users",
		userReactions: []models.UserReaction{
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
		checks: []reactionsGET_check{
			{
				fake.NamespaceId,
				fake.EntityId,
				fake.UserId1,
				&models.ReactionsGETResponse200{
					EntityId: fake.EntityId,
					ReactionsCount: []models.ReactionCount{
						{
							ReactionId: fake.ReactionId,
							Count:      2,
						},
					},
					UserReactions: models.UserReactionsWithinEntity{
						UserId:    fake.UserId1,
						Reactions: []string{fake.ReactionId},
					},
				},
			},
			{
				fake.NamespaceId,
				fake.EntityId,
				fake.UserId2,
				&models.ReactionsGETResponse200{
					EntityId: fake.EntityId,
					ReactionsCount: []models.ReactionCount{
						{
							ReactionId: fake.ReactionId,
							Count:      2,
						},
					},
					UserReactions: models.UserReactionsWithinEntity{
						UserId:    fake.UserId2,
						Reactions: []string{fake.ReactionId},
					},
				},
			},
		},
	},
	{
		name: "different reactions from other users",
		userReactions: []models.UserReaction{
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId1,
				ReactionId:  fake.FreeReactionId1,
			},
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId2,
				ReactionId:  fake.FreeReactionId1,
			},
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId2,
				ReactionId:  fake.FreeReactionId2,
			},
			{
				NamespaceId: fake.NamespaceId,
				EntityId:    fake.EntityId,
				UserId:      fake.UserId3,
				ReactionId:  fake.FreeReactionId2,
			},
		},
		checks: []reactionsGET_check{
			{
				fake.NamespaceId,
				fake.EntityId,
				fake.UserId1,
				&models.ReactionsGETResponse200{
					EntityId: fake.EntityId,
					ReactionsCount: []models.ReactionCount{
						{
							ReactionId: fake.FreeReactionId1,
							Count:      2,
						},
						{
							ReactionId: fake.FreeReactionId2,
							Count:      2,
						},
					},
					UserReactions: models.UserReactionsWithinEntity{
						UserId:    fake.UserId1,
						Reactions: []string{fake.FreeReactionId1},
					},
				},
			},
			{
				fake.NamespaceId,
				fake.EntityId,
				fake.UserId2,
				&models.ReactionsGETResponse200{
					EntityId: fake.EntityId,
					ReactionsCount: []models.ReactionCount{
						{
							ReactionId: fake.FreeReactionId1,
							Count:      2,
						},
						{
							ReactionId: fake.FreeReactionId2,
							Count:      2,
						},
					},
					UserReactions: models.UserReactionsWithinEntity{
						UserId:    fake.UserId2,
						Reactions: []string{fake.FreeReactionId1, fake.FreeReactionId2},
					},
				},
			},
			{
				fake.NamespaceId,
				fake.EntityId,
				fake.UserId3,
				&models.ReactionsGETResponse200{
					EntityId: fake.EntityId,
					ReactionsCount: []models.ReactionCount{
						{
							ReactionId: fake.FreeReactionId1,
							Count:      2,
						},
						{
							ReactionId: fake.FreeReactionId2,
							Count:      2,
						},
					},
					UserReactions: models.UserReactionsWithinEntity{
						UserId:    fake.UserId3,
						Reactions: []string{fake.FreeReactionId2},
					},
				},
			},
		},
	},
}
