package services

import (
	"context"

	"go.uber.org/fx"

	"github.com/Deimvis/reactionsstorage/src/metrics"
	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/src/storages"
	"github.com/Deimvis/reactionsstorage/src/utils"
)

func NewReactionsService(lc fx.Lifecycle, cs *storages.ConfigurationStorage, rs *storages.ReactionsStorage) *ReactionsService {
	return &ReactionsService{cs: cs, rs: rs}
}

type ReactionsService struct {
	cs *storages.ConfigurationStorage
	rs *storages.ReactionsStorage
}

func (s *ReactionsService) GetUserReactions(ctx context.Context, req models.ReactionsGETRequest) models.Response {
	metrics.Record(func() {
		ctx = storages.CtxAcquireConn(ctx, s.rs)
	}, metrics.GETReactions_Acquire)
	defer storages.CtxReleaseConn(&ctx)

	var reactionsCount map[string]int
	var userUniqReactions []string

	metrics.Record(func() {
		reactionsCount = utils.Must(s.rs.GETReactions_GetEntityReactionsCount(ctx, req.Query.NamespaceId, req.Query.EntityId))
	}, metrics.GETReactions_GetEntityReactionsCount)

	metrics.Record(func() {
		userUniqReactions = utils.Must(s.rs.GETReactions_GetUniqEntityUserReactions(ctx, req.Query.NamespaceId, req.Query.EntityId, req.Query.UserId))
	}, metrics.GETReactions_GetUniqEntityUserReactions)

	resp := models.ReactionsGETResponse200{
		EntityId:       req.Query.EntityId,
		ReactionsCount: models.ReactionCount{}.FromMap(reactionsCount),
		UserReactions: models.UserReactionsWithinEntity{
			UserId:    req.Query.UserId,
			Reactions: userUniqReactions,
		},
	}
	return &resp
}

func (s *ReactionsService) AddUserReaction(ctx context.Context, req models.ReactionsPOSTRequest) models.Response {
	metrics.Record(func() {
		ctx = storages.CtxAcquireConn(ctx, s.rs)
	}, metrics.POSTReactions_Acquire)
	defer storages.CtxReleaseConn(&ctx)

	var namespace *models.Namespace
	var err error
	
	metrics.Record(func() {
		namespace, err = s.cs.GetNamespace(ctx, req.Body.NamespaceId)
	}, metrics.POSTReactions_GetNamespace)
	if err != nil {
		return &models.ReactionsPOSTResponse403{Error: err.Error()}
	}
	
	metrics.Record(func() {
		err = s.rs.AddUserReaction(ctx, req.Body, namespace.MaxUniqReactions, namespace.MutuallyExclusiveReactions, *req.Query.Force)
	}, metrics.POSTReactions_AddUserReaction)
	if err != nil {
		return &models.ReactionsPOSTResponse403{Error: err.Error()}
	}

	return &models.ReactionsPOSTResponse200{Status: "ok"}
}

func (s *ReactionsService) RemoveUserReaction(ctx context.Context, req models.ReactionsDELETERequest) models.Response {
	metrics.Record(func() {
		ctx = storages.CtxAcquireConn(ctx, s.rs)
	}, metrics.DELETEReactions_Acquire)
	defer storages.CtxReleaseConn(&ctx)

	var err error
	metrics.Record(func() {
		err = s.rs.RemoveUserReaction(ctx, req.Body)
	}, metrics.DELETEReactions_RemoveUserReaction)
	if err != nil {
		return &models.ReactionsDELETEResponse403{Error: err.Error()}
	}
	return &models.ReactionsDELETEResponse200{Status: "ok"}
}
