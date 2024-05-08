package storages

import (
	"context"
	"errors"

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
	rIds := []string{}
	rows, err := pg.Query(ctx, sql.GetUniqueEntityUserReactions, namespaceId, entityId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var reactionId string
		err := rows.Scan(&reactionId)
		if err != nil {
			return nil, err
		}
		rIds = append(rIds, reactionId)
	}
	return rIds, nil
}

// getEntityReactionsCount returns only reactions with positive count (reactiosn with zero count can be stored physically)
func (rs *ReactionsStorage) getEntityReactionsCount(pg PG, ctx context.Context, namespaceId string, entityId string) (map[string]int, error) {
	rows, err := pg.Query(ctx, sql.GetEntityReactionsCount, namespaceId, entityId)
	if errors.Is(err, pgx.ErrNoRows) {
		err = nil
	}
	defer rows.Close()
	res := make(map[string]int)
	for rows.Next() {
		var reactionId string
		var cnt int
		err := rows.Scan(&reactionId, &cnt)
		if err != nil {
			return nil, err
		}
		res[reactionId] = cnt
	}
	utils.FilterMapIn(res, func(rId string, cnt int) bool { return res[rId] > 0 })
	return res, nil
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
	if err != nil {
		return err
	}
	return nil
}

func (rs *ReactionsStorage) removeConflictingReactions(pg PG, ctx context.Context, newReaction models.UserReaction, reactionsCount map[string]int, uniqUserReactions *[]string, mutExclReactions [][]string) {
	for _, rId := range getConflictingReactionIds(newReaction.ReactionId, *uniqUserReactions, mutExclReactions) {
		// remove from database
		r := models.UserReaction{
			NamespaceId: newReaction.NamespaceId,
			EntityId:    newReaction.EntityId,
			UserId:      newReaction.UserId,
			ReactionId:  rId,
		}
		rs.removeUserReaction(pg, ctx, r)

		// remove from local objects
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
