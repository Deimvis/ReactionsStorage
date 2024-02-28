package storages

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/src/sql"
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
	_, err := rs.pool.Exec(ctx, sql.InitReactionsStorage)
	return err
}

func (rs *ReactionsStorage) GetMutuallyExclusiveReactions(namespaceId string) ([][]string, error) {
	var result [][]string
	err := rs.pool.QueryRow(context.Background(), sql.GetMutuallyExclusiveReactions, namespaceId).Scan(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (rs *ReactionsStorage) GetMutuallyExclusiveReactionsStrict(namespaceId string) [][]string {
	res, err := rs.GetMutuallyExclusiveReactions(namespaceId)
	if err != nil {
		log.Panicf("failed to get mutually exclusive reactions: %s", err)
	}
	return res
}

func (rs *ReactionsStorage) GetMaxUniqueReactions(namespaceId string) (int, error) {
	var result int
	err := rs.pool.QueryRow(context.Background(), sql.GetMaxUniqueReactions, namespaceId).Scan(&result)
	if err != nil {
		return -1, err
	}
	return result, nil
}

func (rs *ReactionsStorage) GetMaxUniqueReactionsStrict(namespaceId string) int {
	res, err := rs.GetMaxUniqueReactions(namespaceId)
	if err != nil {
		log.Panicf("failed to get max unique reactions: %s", err)
	}
	return res
}

func (rs *ReactionsStorage) AddUserReaction(ctx context.Context, reaction models.UserReaction, maxUniqReactions int, mutExclReactions [][]string) error {
	fmt.Println("AddUserReaction")
	tx := rs.beginTxStrict(ctx)
	defer tx.Rollback(ctx)
	lockKey := fmt.Sprintf("%s/%s", reaction.NamespaceId, reaction.EntityId)
	rs.advLockStrict(ctx, tx, lockKey)
	// TODO: check idempotency

	err := rs.addUserReaction(ctx, tx, reaction, maxUniqReactions, mutExclReactions)
	if err != nil {
		return err
	}

	rs.advUnlockStrict(ctx, tx, lockKey)
	tx.Commit(ctx)
	return nil
}

func (rs *ReactionsStorage) GetEntityReactionsCount(ctx context.Context, namespaceId string, entityId string) ([]models.ReactionCount, error) {
	rows, err := rs.pool.Query(ctx, sql.GetEntityReactionsCount, namespaceId, entityId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[models.ReactionCount])
}

func (rs *ReactionsStorage) GetEntityReactionsCountStrict(ctx context.Context, namespaceId string, entityId string) []models.ReactionCount {
	res, err := rs.GetEntityReactionsCount(ctx, namespaceId, entityId)
	if err != nil {
		log.Panicf("failed to get user reactions count for entity: %s", err)
	}
	return res
}

func (rs *ReactionsStorage) addUserReaction(ctx context.Context, tx pgx.Tx, reaction models.UserReaction, maxUniqReactions int, mutExclReactions [][]string) error {
	uniqEntityReactions := rs.getUniqEntityReactionsStrict(ctx, reaction.NamespaceId, reaction.EntityId)
	uniqEntityUserReactions := rs.GetUniqEntityUserReactionsStrict(ctx, reaction.NamespaceId, reaction.EntityId, reaction.UserId)
	err := checkAddUserReaction(ctx, reaction.UserId, reaction.ReactionId, uniqEntityReactions, uniqEntityUserReactions, maxUniqReactions, mutExclReactions)
	// TODO: remove conflicting reactions on `force` flag
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, sql.AddUserReaction, reaction.NamespaceId, reaction.EntityId, reaction.ReactionId, reaction.UserId)
	if err != nil {
		return err
	}
	return nil
}

func (rs *ReactionsStorage) getUniqEntityReactions(ctx context.Context, namespaceId string, entityId string) (map[string]struct{}, error) {
	rows, err := rs.pool.Query(ctx, sql.GetUniqueEntityReactions, namespaceId, entityId)
	if err != nil {
		return nil, err
	}
	return scanUniqReactions(ctx, rows)
}

func (rs *ReactionsStorage) getUniqEntityReactionsStrict(ctx context.Context, namespaceId string, entityId string) map[string]struct{} {
	res, err := rs.getUniqEntityReactions(ctx, namespaceId, entityId)
	if err != nil {
		log.Panicf("failed to get unique entity reactions: %s", err)
	}
	return res
}

func (rs *ReactionsStorage) GetUniqEntityUserReactions(ctx context.Context, namespaceId string, entityId string, userId string) (map[string]struct{}, error) {
	rows, err := rs.pool.Query(ctx, sql.GetUniqueEntityUserReactions, namespaceId, entityId, userId)
	if err != nil {
		return nil, err
	}
	return scanUniqReactions(ctx, rows)
}

func (rs *ReactionsStorage) GetUniqEntityUserReactionsStrict(ctx context.Context, namespaceId string, entityId string, userId string) map[string]struct{} {
	res, err := rs.GetUniqEntityUserReactions(ctx, namespaceId, entityId, userId)
	if err != nil {
		log.Panicf("failed to get unique entity user reactions: %s", err)
	}
	return res
}

func (rs *ReactionsStorage) beginTx(ctx context.Context) (pgx.Tx, error) {
	return rs.pool.Begin(ctx)
}

func (rs *ReactionsStorage) beginTxStrict(ctx context.Context) pgx.Tx {
	tx, err := rs.beginTx(ctx)
	if err != nil {
		log.Panicf("failed to create pg transaction: %s", err)
	}
	return tx
}

func (rs *ReactionsStorage) advLock(ctx context.Context, tx pgx.Tx, key string) error {
	_, err := tx.Exec(ctx, "SELECT pg_advisory_lock( hashtext($1) )", key)
	return err
}

func (rs *ReactionsStorage) advLockStrict(ctx context.Context, tx pgx.Tx, key string) {
	err := rs.advLock(ctx, tx, key)
	if err != nil {
		log.Panicf("failed to acquire advisory lock: %s", err)
	}
}

func (rs *ReactionsStorage) advUnlock(ctx context.Context, tx pgx.Tx, key string) error {
	_, err := tx.Exec(ctx, "SELECT pg_advisory_unlock( hashtext($1) )", key)
	return err
}

func (rs *ReactionsStorage) advUnlockStrict(ctx context.Context, tx pgx.Tx, key string) {
	err := rs.advUnlock(ctx, tx, key)
	if err != nil {
		log.Panicf("failed to release advisory lock: %s", err)
	}
}
