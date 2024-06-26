package storages

import (
	"context"
	"errors"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Deimvis/reactionsstorage/src/models"
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

func (cs *ConfigurationStorage) GetPool() *pgxpool.Pool {
	return cs.pool
}

func (cs *ConfigurationStorage) Init(ctx context.Context) error {
	return cs.init(AcquirePG(ctx, cs), ctx)
}

func (cs *ConfigurationStorage) AddReaction(ctx context.Context, r *models.Reaction) error {
	return cs.addReaction(AcquirePG(ctx, cs), ctx, r)
}

func (cs *ConfigurationStorage) AddReactionSet(ctx context.Context, r *models.ReactionSet) error {
	return cs.addReactionSet(AcquirePG(ctx, cs), ctx, r)
}

func (cs *ConfigurationStorage) AddNamespace(ctx context.Context, n *models.Namespace) error {
	return cs.addNamespace(AcquirePG(ctx, cs), ctx, n)
}

func (cs *ConfigurationStorage) GetNamespace(ctx context.Context, namespaceId string) (*models.Namespace, error) {
	return cs.getNamespace(AcquirePG(ctx, cs), ctx, namespaceId)
}

func (cs *ConfigurationStorage) HasNamespace(ctx context.Context, namespaceId string) bool {
	_, err := cs.GetNamespace(ctx, namespaceId)
	return !errors.Is(err, pgx.ErrNoRows)
}

func (cs *ConfigurationStorage) GetAvailableReactions(ctx context.Context, namespaceId string) ([]models.Reaction, error) {
	return cs.getAvailableReactions(AcquirePG(ctx, cs), ctx, namespaceId)
}

func (cs *ConfigurationStorage) GetMutuallyExclusiveReactions(ctx context.Context, namespaceId string) ([][]string, error) {
	return cs.getMutuallyExclusiveReactions(AcquirePG(ctx, cs), ctx, namespaceId)
}

func (cs *ConfigurationStorage) GetMaxUniqueReactions(ctx context.Context, namespaceId string) (int, error) {
	return cs.getMaxUniqueReactions(AcquirePG(ctx, cs), ctx, namespaceId)
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

func (cs *ConfigurationStorage) Clear(ctx context.Context) error {
	return cs.clear(AcquirePG(ctx, cs), ctx)
}

func (cs *ConfigurationStorage) beginTx(ctx context.Context) (pgx.Tx, error) {
	return AcquirePG(ctx, cs).Begin(ctx)
}
