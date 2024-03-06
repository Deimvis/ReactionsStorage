package storages

import (
	"context"
	"fmt"

	"github.com/Deimvis/reactionsstorage/src/models"
)

func (cs *ConfigurationStorage) AddReactionStrict(ctx context.Context, r *models.Reaction) {
	err := cs.AddReaction(ctx, r)
	if err != nil {
		panic(fmt.Errorf("failed to add reaction: %w", err))
	}
}

func (cs *ConfigurationStorage) AddReactionSetStrict(ctx context.Context, r *models.ReactionSet) {
	err := cs.AddReactionSet(ctx, r)
	if err != nil {
		panic(fmt.Errorf("failed to add reaction set: %w", err))
	}
}

func (cs *ConfigurationStorage) AddNamespaceStrict(ctx context.Context, n *models.Namespace) {
	err := cs.AddNamespace(ctx, n)
	if err != nil {
		panic(fmt.Errorf("failed to add namespace: %w", err))
	}
}

func (cs *ConfigurationStorage) GetAvailableReactionsStrict(ctx context.Context, namespaceId string) []models.Reaction {
	res, err := cs.GetAvailableReactions(ctx, namespaceId)
	fmt.Println("res", res)
	fmt.Println(err)
	if err != nil {
		panic(fmt.Errorf("failed to get avaialble reactions: %w", err))
	}
	return res
}

func (cs *ConfigurationStorage) GetMutuallyExclusiveReactionsStrict(namespaceId string) [][]string {
	res, err := cs.GetMutuallyExclusiveReactions(namespaceId)
	if err != nil {
		panic(fmt.Errorf("failed to get mutually exclusive reactions: %w", err))
	}
	return res
}

func (cs *ConfigurationStorage) GetMaxUniqueReactionsStrict(namespaceId string) int {
	res, err := cs.GetMaxUniqueReactions(namespaceId)
	if err != nil {
		panic(fmt.Errorf("failed to get max unique reactions: %w", err))
	}
	return res
}

func (cs *ConfigurationStorage) ClearStrict(ctx context.Context) {
	err := cs.Clear(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to clear configuration storage: %w", err))
	}
}
