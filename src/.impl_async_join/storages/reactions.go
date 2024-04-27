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

func (cs *ReactionsStorage) GetPool() *pgxpool.Pool {
	return cs.pool
}

func (rs *ReactionsStorage) Init(ctx context.Context) error {
	return rs.init(AcquirePG(ctx, rs), ctx)
}

// GetEntityReactionsCount erturns only reactions with positive count (reactiosn with zero count can be stored physically)
func (rs *ReactionsStorage) GetEntityReactionsCount(ctx context.Context, namespaceId string, entityId string) (map[string]int, error) {
	return rs.getEntityReactionsCount(AcquirePG(ctx, rs), ctx, namespaceId, entityId)
}

func (rs *ReactionsStorage) GetUniqEntityUserReactions(ctx context.Context, namespaceId string, entityId string, userId string) ([]string, error) {
	return rs.getUniqEntityUserReactions(AcquirePG(ctx, rs), ctx, namespaceId, entityId, userId)
}

func (rs *ReactionsStorage) AddUserReaction(ctx context.Context, reaction models.UserReaction, maxUniqReactions int, mutExclReactions [][]string, force bool) error {
	tx := utils.Must(rs.beginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}))
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
	return rs.removeUserReaction(AcquirePG(ctx, rs), ctx, reaction)
}

// GetUserReactions is supposed to be used only for debug and test purposes
func (rs *ReactionsStorage) GetUserReactions(ctx context.Context) ([]models.UserReaction, error) {
	return rs.getUserReactions(AcquirePG(ctx, rs), ctx)
}

func (rs *ReactionsStorage) Clear(ctx context.Context) error {
	return rs.clear(AcquirePG(ctx, rs), ctx)
}

func (rs *ReactionsStorage) RefreshEntityReactions(ctx context.Context) error {
	return rs.refreshEntityReactions(AcquirePG(ctx, rs), ctx)
}

func (rs *ReactionsStorage) beginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return AcquirePG(ctx, rs).BeginTx(ctx, txOptions)
}

func (rs *ReactionsStorage) advLock(ctx context.Context, tx pgx.Tx, key string) error {
	_, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock( hashtext($1) );", key)
	return err
}

func (rs *ReactionsStorage) advUnlock(ctx context.Context, tx pgx.Tx, key string) error {
	_, err := tx.Exec(ctx, "SELECT pg_advisory_xact_unlock( hashtext($1) )", key)
	return err
}
