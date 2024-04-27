package storages

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/src/sql"
	"github.com/Deimvis/reactionsstorage/src/utils"
)

func (rs *ReactionsStorage) init(pg PG, ctx context.Context) error {
	_, err := pg.Exec(ctx, sql.InitReactionsStorage)
	return err
}

func (rs *ReactionsStorage) getUniqEntityUserReactions(pg PG, ctx context.Context, namespaceId string, entityId string, userId string) ([]string, error) {
	res := []string{}
	rows, err := pg.Query(ctx, sql.GetUniqueEntityUserReactions, namespaceId, entityId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var rId string
		err := rows.Scan(&rId)
		if err != nil {
			return nil, err
		}
		res = append(res, rId)
	}
	return res, nil
}

// getEntityReactionsCount erturns only reactions with positive count (reactiosn with zero count can be stored physically)
func (rs *ReactionsStorage) getEntityReactionsCount(pg PG, ctx context.Context, namespaceId string, entityId string) (map[string]int, error) {
	res := make(map[string]int)
	var err error
	key := fmt.Sprintf("%s__%s", namespaceId, entityId)
	err = pg.QueryRow(ctx, sql.GetEntityReactionsCount, key).Scan(&res)
	if errors.Is(err, pgx.ErrNoRows) {
		err = nil
	}
	return res, err
}

func (rs *ReactionsStorage) addUserReaction(pg PG, ctx context.Context, reaction models.UserReaction, maxUniqReactions int, mutExclReactions [][]string, force bool) error {
	reactionsCount := utils.Must(rs.getEntityReactionsCount(pg, ctx, reaction.NamespaceId, reaction.EntityId))
	uniqUserReactions := utils.Must(rs.getUniqEntityUserReactions(pg, ctx, reaction.NamespaceId, reaction.EntityId, reaction.UserId))

	err := checkAddUserReaction(ctx, reaction.UserId, reaction.ReactionId, reactionsCount, uniqUserReactions, maxUniqReactions, mutExclReactions)
	for err != nil {
		_, isConflict := err.(*ConflictingReactionError)
		if isConflict && force {
			rs.removeConflictingReactions(pg, ctx, reaction, reactionsCount, &uniqUserReactions, mutExclReactions)
			force = false
		} else {
			break
		}
		err = checkAddUserReaction(ctx, reaction.UserId, reaction.ReactionId, reactionsCount, uniqUserReactions, maxUniqReactions, mutExclReactions)
	}
	if err != nil {
		return err
	}

	_, err = pg.Exec(ctx, sql.AddUserReaction, reaction.NamespaceId, reaction.EntityId, reaction.ReactionId, reaction.UserId)
	return err
}

func (rs *ReactionsStorage) removeConflictingReactions(pg PG, ctx context.Context, newReaction models.UserReaction, reactionsCount map[string]int, uniqUserReactions *[]string, mutExclReactions [][]string) {
	for _, rId := range getConflictingReactionIds(newReaction.ReactionId, *uniqUserReactions, mutExclReactions) {
		r := models.UserReaction{
			NamespaceId: newReaction.NamespaceId,
			EntityId:    newReaction.EntityId,
			UserId:      newReaction.UserId,
			ReactionId:  rId,
		}
		rs.removeUserReaction(pg, ctx, r)
		utils.FilterIn(uniqUserReactions, func(el string) bool { return el != rId })
		reactionsCount[rId] -= 1
		if reactionsCount[rId] == 0 {
			delete(reactionsCount, rId)
		}
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

func (rs *ReactionsStorage) clear(pg PG, ctx context.Context) error {
	_, err := pg.Exec(ctx, sql.ClearUserReactionsStorage)
	return err
}

func (rs *ReactionsStorage) refreshEntityReactions(pg PG, ctx context.Context) error {
	_, err := pg.Exec(ctx, sql.RefreshEntityReactions)
	return err
}

// execBatch returns first error occured
func execBatch(pg PG, ctx context.Context, batch *pgx.Batch) error {
	results := pg.SendBatch(ctx, batch)
	defer results.Close()

	for i := 0; i < batch.Len(); i++ {
		_, err := results.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}
