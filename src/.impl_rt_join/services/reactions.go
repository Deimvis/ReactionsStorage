package services

import (
	"context"

	"go.uber.org/fx"

	"github.com/Deimvis/reactionsstorage/src/metrics"
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

func (s *ReactionsService) GetUserReactions(ctx context.Context, req models.ReactionsGETRequest) models.Response {
	// TODO: remove debug metric wrappers

	// ctx = storages.CtxAcquireConn(ctx, s.rs)
	metrics.Record(func() {
		ctx = storages.CtxAcquireConn(ctx, s.rs)
	}, metrics.GETReactionsAcquire)
	defer storages.CtxReleaseConn(&ctx)

	var reactionsCount []models.ReactionCount
	var userUniqReactions map[string]struct{}

	// reactionsCount = s.rs.GetEntityReactionsCountStrict(ctx, req.Query.NamespaceId, req.Query.EntityId)
	metrics.Record(func() {
		reactionsCount = s.rs.GetEntityReactionsCountStrict(ctx, req.Query.NamespaceId, req.Query.EntityId)
	}, metrics.GetEntityReactionsCount)
	// userUniqReactions = s.rs.GetUniqEntityUserReactionsStrict(ctx, req.Query.NamespaceId, req.Query.EntityId, req.Query.UserId)
	metrics.Record(func() {
		userUniqReactions = s.rs.GetUniqEntityUserReactionsStrict(ctx, req.Query.NamespaceId, req.Query.EntityId, req.Query.UserId)
	}, metrics.GetUniqEntityUserReactions)

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

func (s *ReactionsService) AddUserReaction(ctx context.Context, req models.ReactionsPOSTRequest) models.Response {
	ctx = storages.CtxAcquireConn(ctx, s.rs)
	defer storages.CtxReleaseConn(&ctx)

	namespace, err := s.cs.GetNamespace(ctx, req.Body.NamespaceId)
	if err != nil {
		return &models.ReactionsPOSTResponse403{Error: err.Error()}
	}
	err = s.rs.AddUserReaction(ctx, req.Body, namespace.MaxUniqReactions, namespace.MutuallyExclusiveReactions, *req.Query.Force)
	if err != nil {
		return &models.ReactionsPOSTResponse403{Error: err.Error()}
	}
	return &models.ReactionsPOSTResponse200{Status: "ok"}
}

func (s *ReactionsService) RemoveUserReaction(ctx context.Context, req models.ReactionsDELETERequest) models.Response {
	err := s.rs.RemoveUserReaction(ctx, req.Body)
	if err != nil {
		return &models.ReactionsDELETEResponse403{Error: err.Error()}
	}
	return &models.ReactionsDELETEResponse200{Status: "ok"}
}