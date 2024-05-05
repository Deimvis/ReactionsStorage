package storages

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	err := pg.QueryRow(ctx, sql.GetUniqueEntityUserReactions, namespaceId, entityId, userId).Scan(&rIds)
	if errors.Is(err, pgx.ErrNoRows) {
		err = nil
	}
	return rIds, err
}

// getEntityReactionsCount erturns only reactions with positive count (reactiosn with zero count can be stored physically)
func (rs *ReactionsStorage) getEntityReactionsCount(pg PG, ctx context.Context, namespaceId string, entityId string) (map[string]int, error) {
	res := make(map[string]int)
	err := pg.QueryRow(ctx, sql.GETReactions_GetEntityReactionsCount, namespaceId, entityId).Scan(&res)
	if errors.Is(err, pgx.ErrNoRows) {
		err = nil
	}
	utils.FilterMapIn(res, func(rId string, cnt int) bool { return res[rId] > 0 })
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

	queries := sql.ParseQueries(sql.AddUserReaction)
	if len(queries) != 2 {
		panic(fmt.Errorf("add_user_reaction query has wrong number of queries: %d (2 was expected)", len(queries)))
	}
	batch := &pgx.Batch{}
	batch.Queue(queries[0], reaction.NamespaceId, reaction.EntityId, reaction.ReactionId, reaction.UserId, time.Now().Unix())
	batch.Queue(queries[1], reaction.NamespaceId, reaction.EntityId, reaction.ReactionId)
	fmt.Println("exec adding user reaction query", reactionsCount, uniqUserReactions, reaction.ReactionId)
	return execBatch(pg, ctx, batch)
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
	queries := sql.ParseQueries(sql.RemoveUserReaction)
	if len(queries) != 2 {
		panic(fmt.Errorf("remove_user_reaction query has wrong number of queries: %d (2 was expected)", len(queries)))
	}
	batch := &pgx.Batch{}
	batch.Queue(queries[0], reaction.NamespaceId, reaction.EntityId, reaction.ReactionId, reaction.UserId)
	batch.Queue(queries[1], reaction.NamespaceId, reaction.EntityId, reaction.ReactionId)
	return execBatch(pg, ctx, batch)
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
