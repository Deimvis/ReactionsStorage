package services

import (
	"context"
	"log"

	"go.uber.org/fx"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/src/storages"
)

func NewReactionsService(lc fx.Lifecycle, cs *storages.ConfigurationStorage, rs *storages.ReactionsStorage) *ReactionsService {
	return &ReactionsService{cs: cs, rs: rs}
}

type ReactionsService struct {
	cs *storages.ConfigurationStorage
	rs *storages.ReactionsStorage
}

func (rs *ReactionsService) GetUserReactions(ctx context.Context, req models.ReactionsGETRequest) models.Response {
	reactionsCount := rs.rs.GetEntityReactionsCountStrict(ctx, req.Query.NamespaceId, req.Query.EntityId)
	userUniqReactions := rs.rs.GetUniqEntityUserReactionsStrict(ctx, req.Query.NamespaceId, req.Query.EntityId, req.Query.UserId)
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
	namespace, err := rs.cs.GetNamespace(req.Body.NamespaceId)
	if err != nil {
		return &models.ReactionsPOSTResponse403{Error: err.Error()}
	}
	log.Println("Namespace:", namespace)
	err = rs.rs.AddUserReaction(ctx, req.Body, namespace.MaxUniqReactions, namespace.MutuallyExclusiveReactions)
	if err != nil {
		return &models.ReactionsPOSTResponse403{Error: err.Error()}
	}
	return &models.ReactionsPOSTResponse200{Status: "ok"}
}

func (rs *ReactionsService) RemoveUserReaction(ctx context.Context, req models.ReactionsDELETERequest) models.Response {
	err := rs.rs.RemoveUserReaction(ctx, req.Body)
	if err != nil {
		return &models.ReactionsDELETEResponse403{Error: err.Error()}
	}
	return &models.ReactionsDELETEResponse200{Status: "ok"}
}
