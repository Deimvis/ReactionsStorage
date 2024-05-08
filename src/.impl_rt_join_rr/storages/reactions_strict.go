package storages

import (
	"context"
	"fmt"
)

// func (rs *ReactionsStorage) GetEntityReactionsCountStrict(ctx context.Context, namespaceId string, entityId string) []models.ReactionCount {
// 	res, err := rs.GetEntityReactionsCount(ctx, namespaceId, entityId)
// 	if err != nil {
// 		panic(fmt.Errorf("failed to get user reactions count for entity: %w", err))
// 	}
// 	return res
// }

func (rs *ReactionsStorage) ClearStrict(ctx context.Context) {
	err := rs.Clear(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to clear user reactions: %w", err))
	}
}

// func (rs *ReactionsStorage) GetUniqEntityUserReactionsStrict(ctx context.Context, namespaceId string, entityId string, userId string) map[string]struct{} {
// 	res, err := rs.GetUniqEntityUserReactions(ctx, namespaceId, entityId, userId)
// 	if err != nil {
// 		panic(fmt.Errorf("failed to get unique entity user reactions: %w", err))
// 	}
// 	return res
// }
