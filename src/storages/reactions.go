package storages

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/src/utils"
)

func NewReactionsStorage(lc fx.Lifecycle, pool *pgxpool.Pool) *ReactionsStorage {
	storage := &ReactionsStorage{pool: pool}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return storage.Init(ctx)
		},
	})
	return storage
}

type ReactionsStorage struct {
	pool *pgxpool.Pool
}

func (rs *ReactionsStorage) Init(ctx context.Context) error {
	return rs.init(rs.pool, ctx)
}

func (rs *ReactionsStorage) GetEntityReactionsCount(ctx context.Context, namespaceId string, entityId string) ([]models.ReactionCount, error) {
	return rs.getEntityReactionsCount(rs.pool, ctx, namespaceId, entityId)
}

func (rs *ReactionsStorage) GetUniqEntityUserReactions(ctx context.Context, namespaceId string, entityId string, userId string) (map[string]struct{}, error) {
	return rs.getUniqEntityUserReactions(rs.pool, ctx, namespaceId, entityId, userId)
}

func (rs *ReactionsStorage) AddUserReaction(ctx context.Context, reaction models.UserReaction, maxUniqReactions int, mutExclReactions [][]string, force bool) error {
	tx := utils.Must(rs.beginTx(ctx))
	defer tx.Rollback(ctx)
	lockKey := fmt.Sprintf("%s/%s", reaction.NamespaceId, reaction.EntityId)
	utils.Must0(rs.advLock(ctx, tx, lockKey)) // transaction-level lock is automatically released at the end of tx

	err := rs.addUserReaction(tx, ctx, reaction, maxUniqReactions, mutExclReactions, force)
	if err != nil {
		return err
	}

	tx.Commit(ctx)
	return nil
}

func (rs *ReactionsStorage) RemoveUserReaction(ctx context.Context, reaction models.UserReaction) error {
	return rs.removeUserReaction(rs.pool, ctx, reaction)
}

// GetUserReactions is supposed to be used only for debug and test purposes
func (rs *ReactionsStorage) GetUserReactions(ctx context.Context) ([]models.UserReaction, error) {
	return rs.getUserReactions(rs.pool, ctx)
}

func (rs *ReactionsStorage) Clear(ctx context.Context) error {
	return rs.clear(rs.pool, ctx)
}

func (rs *ReactionsStorage) beginTx(ctx context.Context) (pgx.Tx, error) {
	return rs.pool.Begin(ctx)
}

func (rs *ReactionsStorage) advLock(ctx context.Context, tx pgx.Tx, key string) error {
	_, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock( hashtext($1) );", key)
	return err
}

func (rs *ReactionsStorage) advUnlock(ctx context.Context, tx pgx.Tx, key string) error {
	_, err := tx.Exec(ctx, "SELECT pg_advisory_xact_unlock( hashtext($1) )", key)
	return err
}
