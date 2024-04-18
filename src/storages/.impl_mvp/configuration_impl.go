package storages

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/src/sql"
)


func (cs *ConfigurationStorage) init(pg PG, ctx context.Context) error {
	_, err := pg.Exec(ctx, sql.InitConfigurationStorage)
	return err
}

func (cs *ConfigurationStorage) addReaction(pg PG, ctx context.Context, r *models.Reaction) error {
	_, err := pg.Exec(ctx, sql.AddReaction, r.Id, r.ShortName, r.Type, r.Code, r.Url)
	return err
}

func (cs *ConfigurationStorage) addReactionSet(pg PG, ctx context.Context, r *models.ReactionSet) error {
	_, err := pg.Exec(ctx, sql.AddReactionSet, r.Id, r.ReactionIds)
	return err
}

func (cs *ConfigurationStorage) addNamespace(pg PG, ctx context.Context, n *models.Namespace) error {
	_, err := pg.Exec(ctx, sql.AddNamespace, n.Id, n.ReactionSetId, n.MaxUniqReactions, n.MutuallyExclusiveReactions)
	return err
}

func (cs *ConfigurationStorage) getNamespace(pg PG, ctx context.Context, namespaceId string) (*models.Namespace, error) {
	row, err := pg.Query(ctx, sql.GetNamespace, namespaceId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectExactlyOneRow(row, pgx.RowToAddrOfStructByName[models.Namespace])
}


func (cs *ConfigurationStorage) getAvailableReactions(pg PG, ctx context.Context, namespaceId string) ([]models.Reaction, error) {
	rows, err := pg.Query(ctx, sql.GetAvailableReactions, namespaceId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[models.Reaction])
}

func (cs *ConfigurationStorage) getMutuallyExclusiveReactions(pg PG, ctx context.Context, namespaceId string) ([][]string, error) {
	var result [][]string
	err := pg.QueryRow(ctx, sql.GetMutuallyExclusiveReactions, namespaceId).Scan(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (cs *ConfigurationStorage) getMaxUniqueReactions(pg PG, ctx context.Context, namespaceId string) (int, error) {
	var result int
	err := pg.QueryRow(ctx, sql.GetMaxUniqueReactions, namespaceId).Scan(&result)
	if err != nil {
		return -1, err
	}
	return result, nil
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

func (cs *ConfigurationStorage) clear(pg PG, ctx context.Context) error {
	_, err := pg.Exec(ctx, sql.ClearConfigurationStorage)
	return err
}
