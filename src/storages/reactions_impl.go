package storages

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/src/sql"
	"github.com/Deimvis/reactionsstorage/src/utils"
)

func (rs *ReactionsStorage) getUniqEntityReactions(pg PG, ctx context.Context, namespaceId string, entityId string) (map[string]struct{}, error) {
	rows, err := pg.Query(ctx, sql.GetUniqueEntityReactions, namespaceId, entityId)
	if err != nil {
		return nil, err
	}
	return scanUniqReactions(ctx, rows)
}

func (rs *ReactionsStorage) getUniqEntityUserReactions(pg PG, ctx context.Context, namespaceId string, entityId string, userId string) (map[string]struct{}, error) {
	rows, err := pg.Query(ctx, sql.GetUniqueEntityUserReactions, namespaceId, entityId, userId)
	if err != nil {
		return nil, err
	}
	return scanUniqReactions(ctx, rows)
}

func (rs *ReactionsStorage) getEntityReactionsCount(pg PG, ctx context.Context, namespaceId string, entityId string) ([]models.ReactionCount, error) {
	rows, err := pg.Query(ctx, sql.GetEntityReactionsCount, namespaceId, entityId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[models.ReactionCount])
}

func (rs *ReactionsStorage) addUserReaction(pg PG, ctx context.Context, reaction models.UserReaction, maxUniqReactions int, mutExclReactions [][]string, force bool) error {
	uniqEntityReactions := utils.Must(rs.getUniqEntityReactions(pg, ctx, reaction.NamespaceId, reaction.EntityId))
	uniqEntityUserReactions := utils.Must(rs.getUniqEntityUserReactions(pg, ctx, reaction.NamespaceId, reaction.EntityId, reaction.UserId))

	err := checkAddUserReaction(ctx, reaction.UserId, reaction.ReactionId, uniqEntityReactions, uniqEntityUserReactions, maxUniqReactions, mutExclReactions)
	for err != nil {
		_, isConflict := err.(*ConflictingReactionError)
		if isConflict && force {
			rs.removeConflictingReactions(pg, ctx, reaction, uniqEntityUserReactions, mutExclReactions)
			force = false
		} else {
			break
		}
		err = checkAddUserReaction(ctx, reaction.UserId, reaction.ReactionId, uniqEntityReactions, uniqEntityUserReactions, maxUniqReactions, mutExclReactions)
	}
	if err != nil {
		return err
	}

	_, err = pg.Exec(ctx, sql.AddUserReaction, reaction.NamespaceId, reaction.EntityId, reaction.ReactionId, reaction.UserId)
	if err != nil {
		return err
	}
	return nil
}

func (rs *ReactionsStorage) removeConflictingReactions(pg PG, ctx context.Context, newReaction models.UserReaction, uniqEntityUserReactions map[string]struct{}, mutExclReactions [][]string) {
	for _, rId := range getConflictingReactionIds(newReaction.ReactionId, uniqEntityUserReactions, mutExclReactions) {
		r := models.UserReaction{
			NamespaceId: newReaction.NamespaceId,
			EntityId:    newReaction.EntityId,
			UserId:      newReaction.UserId,
			ReactionId:  rId,
		}
		rs.removeUserReaction(pg, ctx, r)
		delete(uniqEntityUserReactions, rId)
	}
}

func (rs *ReactionsStorage) removeUserReaction(pg PG, ctx context.Context, reaction models.UserReaction) error {
	_, err := pg.Exec(ctx, sql.RemoveUserReaction, reaction.NamespaceId, reaction.EntityId, reaction.ReactionId, reaction.UserId)
	return err
}

func (rs *ReactionsStorage) getUserReactions(pg PG, ctx context.Context) ([]models.UserReaction, error) {
	rows, err := pg.Query(ctx, sql.GetUserReactions)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[models.UserReaction])
}
