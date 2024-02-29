package services

import (
	"context"
	"log"

	"go.uber.org/fx"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/src/storages"
)

func NewReactionsService(lc fx.Lifecycle, storage *storages.ReactionsStorage) *ReactionsService {
	return &ReactionsService{storage: storage}
}

type ReactionsService struct {
	storage *storages.ReactionsStorage
}

func (rs *ReactionsService) GetUserReactions(ctx context.Context, req models.ReactionsGETRequest) models.Response {
	reactionsCount := rs.storage.GetEntityReactionsCountStrict(ctx, req.Query.NamespaceId, req.Query.EntityId)
	userUniqReactions := rs.storage.GetUniqEntityUserReactionsStrict(ctx, req.Query.NamespaceId, req.Query.EntityId, req.Query.UserId)
	resp := models.ReactionsGETResponse200{
		EntityId:       req.Query.EntityId,
		ReactionsCount: reactionsCount,
		UserReactions: models.UserReactionsWithinEntity{
			UserId:    req.Query.UserId,
			Reactions: GetKeys(userUniqReactions),
		},
	}
	return &resp
}

func (rs *ReactionsService) AddUserReaction(ctx context.Context, req models.ReactionsPOSTRequest) models.Response {
	maxUniqReactions := rs.storage.GetMaxUniqueReactionsStrict(req.Body.NamespaceId)
	mutExclReactions := rs.storage.GetMutuallyExclusiveReactionsStrict(req.Body.NamespaceId)
	log.Println(maxUniqReactions, mutExclReactions)
	err := rs.storage.AddUserReaction(ctx, req.Body, maxUniqReactions, mutExclReactions)
	if err != nil {
		return &models.ReactionsPOSTResponse403{Error: err.Error()}
	}
	return &models.ReactionsPOSTResponse200{Status: "ok"}
}

func (rs *ReactionsService) RemoveUserReaction(ctx context.Context, req models.ReactionsDELETERequest) models.Response {
	err := rs.storage.RemoveUserReaction(ctx, req.Body)
	if err != nil {
		return &models.ReactionsDELETEResponse403{Error: err.Error()}
	}
	return &models.ReactionsDELETEResponse200{Status: "ok"}
}