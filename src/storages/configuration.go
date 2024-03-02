package storages

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/src/sql"
)

func NewConfigurationStorage(lc fx.Lifecycle, pool *pgxpool.Pool) *ConfigurationStorage {
	storage := &ConfigurationStorage{pool: pool}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return storage.Init(ctx)
		},
	})
	return storage
}

type ConfigurationStorage struct {
	pool *pgxpool.Pool
}

func (cs *ConfigurationStorage) Init(ctx context.Context) error {
	_, err := cs.pool.Exec(ctx, sql.InitConfigurationStorage)
	return err
}

func (cs *ConfigurationStorage) AddReaction(ctx context.Context, r *models.Reaction) error {
	_, err := cs.pool.Exec(ctx, sql.AddReaction, r.Id, r.ShortName, r.Type, r.Code, r.Url)
	return err
}

func (cs *ConfigurationStorage) AddReactionSet(ctx context.Context, r *models.ReactionSet) error {
	_, err := cs.pool.Exec(ctx, sql.AddReactionSet, r.Id, r.ReactionIds)
	return err
}

func (cs *ConfigurationStorage) AddNamespace(ctx context.Context, n *models.Namespace) error {
	_, err := cs.pool.Exec(ctx, sql.AddNamespace, n.Id, n.ReactionSetId, n.MaxUniqReactions, n.MutuallyExclusiveReactions)
	return err
}

func (cs *ConfigurationStorage) GetNamespace(namespaceId string) (*models.Namespace, error) {
	row, err := cs.pool.Query(context.Background(), sql.GetNamespace, namespaceId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectExactlyOneRow(row, pgx.RowToAddrOfStructByName[models.Namespace])
}

func (cs *ConfigurationStorage) GetMutuallyExclusiveReactions(namespaceId string) ([][]string, error) {
	var result [][]string
	err := cs.pool.QueryRow(context.Background(), sql.GetMutuallyExclusiveReactions, namespaceId).Scan(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (cs *ConfigurationStorage) GetMaxUniqueReactions(namespaceId string) (int, error) {
	var result int
	err := cs.pool.QueryRow(context.Background(), sql.GetMaxUniqueReactions, namespaceId).Scan(&result)
	if err != nil {
		return -1, err
	}
	return result, nil
}

func (cs *ConfigurationStorage) Clear(ctx context.Context) error {
	_, err := cs.pool.Exec(ctx, sql.ClearConfigurationStorage)
	return err
}
