package storages

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/src/sql"
	"github.com/Deimvis/reactionsstorage/src/utils"
)

func NewConfigurationStorage(lc fx.Lifecycle, pool *pgxpool.Pool, logger *zap.SugaredLogger) *ConfigurationStorage {
	storage := &ConfigurationStorage{pool: pool, logger: logger}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return storage.Init(ctx)
		},
	})
	return storage
}

type ConfigurationStorage struct {
	pool *pgxpool.Pool

	logger *zap.SugaredLogger
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

func (cs *ConfigurationStorage) GetNamespace(ctx context.Context, namespaceId string) (*models.Namespace, error) {
	row, err := cs.pool.Query(ctx, sql.GetNamespace, namespaceId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectExactlyOneRow(row, pgx.RowToAddrOfStructByName[models.Namespace])
}

func (cs *ConfigurationStorage) HasNamespace(ctx context.Context, namespaceId string) bool {
	_, err := cs.GetNamespace(ctx, namespaceId)
	return !errors.Is(err, pgx.ErrNoRows)
}

func (cs *ConfigurationStorage) GetAvailableReactions(ctx context.Context, namespaceId string) ([]models.Reaction, error) {
	rows, err := cs.pool.Query(ctx, sql.GetAvailableReactions, namespaceId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[models.Reaction])
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

func (cs *ConfigurationStorage) SetConfiguration(ctx context.Context, conf *models.Configuration) error {
	tx := utils.Must(cs.beginTx(ctx))
	defer tx.Rollback(ctx)

	err := utils.UntilFirstErr(
		func() error { return cs.clear(tx, ctx) },
		func() error { return cs.addReactionsBatch(tx, ctx, conf.Reactions) },
		func() error { return cs.addReactionSetsBatch(tx, ctx, conf.ReactionSets) },
		func() error { return cs.addNamespacesBatch(tx, ctx, conf.Namespaces) },
	)
	if err != nil {
		return err
	}

	tx.Commit(ctx)
	return nil
}

func (cs *ConfigurationStorage) addReactionsBatch(pg PG, ctx context.Context, reactions []models.Reaction) error {
	_, err := pg.CopyFrom(
		ctx,
		pgx.Identifier{"reaction"},
		[]string{"id", "short_name", "type", "code", "url"},
		pgx.CopyFromSlice(len(reactions), func(i int) ([]any, error) {
			return []any{
				reactions[i].Id,
				reactions[i].ShortName,
				reactions[i].Type,
				reactions[i].Code,
				reactions[i].Url,
			}, nil
		}),
	)
	cs.logger.Infow("Batch insert into `reaction`", "error", err, "size", len(reactions))
	return err
}

func (cs *ConfigurationStorage) addReactionSetsBatch(pg PG, ctx context.Context, reactionSets []models.ReactionSet) error {
	_, err := pg.CopyFrom(
		ctx,
		pgx.Identifier{"reaction_set"},
		[]string{"id", "reaction_ids"},
		pgx.CopyFromSlice(len(reactionSets), func(i int) ([]any, error) {
			return []any{
				reactionSets[i].Id,
				reactionSets[i].ReactionIds,
			}, nil
		}),
	)
	cs.logger.Infow("Batch insert into `reaction_set`", "error", err, "size", len(reactionSets))
	return err
}

func (cs *ConfigurationStorage) addNamespacesBatch(pg PG, ctx context.Context, namespaces []models.Namespace) error {
	_, err := pg.CopyFrom(
		ctx,
		pgx.Identifier{"namespace"},
		[]string{"id", "reaction_set_id", "max_uniq_reactions", "mutually_exclusive_reactions"},
		pgx.CopyFromSlice(len(namespaces), func(i int) ([]any, error) {
			return []any{
				namespaces[i].Id,
				namespaces[i].ReactionSetId,
				namespaces[i].MaxUniqReactions,
				namespaces[i].MutuallyExclusiveReactions,
			}, nil
		}),
	)
	cs.logger.Infow("Batch insert into `namespaces`", "error", err, "size", len(namespaces))
	return err
}

func (cs *ConfigurationStorage) Clear(ctx context.Context) error {
	return cs.clear(cs.pool, ctx)
}

func (cs *ConfigurationStorage) clear(pg PG, ctx context.Context) error {
	_, err := cs.pool.Exec(ctx, sql.ClearConfigurationStorage)
	return err
}

func (cs *ConfigurationStorage) beginTx(ctx context.Context) (pgx.Tx, error) {
	return cs.pool.Begin(ctx)
}
