package storages

import (
	"context"
	"log"

	"github.com/Deimvis/reactionsstorage/src/models"
)

func (cs *ConfigurationStorage) AddReactionStrict(ctx context.Context, r *models.Reaction) {
	err := cs.AddReaction(ctx, r)
	if err != nil {
		log.Panicf("failed to add reaction: %s", err)
	}
}

func (cs *ConfigurationStorage) AddReactionSetStrict(ctx context.Context, r *models.ReactionSet) {
	err := cs.AddReactionSet(ctx, r)
	if err != nil {
		log.Panicf("failed to add reaction set: %s", err)
	}
}

func (cs *ConfigurationStorage) AddNamespaceStrict(ctx context.Context, n *models.Namespace) {
	err := cs.AddNamespace(ctx, n)
	if err != nil {
		log.Panicf("failed to add namespace: %s", err)
	}
}

func (cs *ConfigurationStorage) GetMutuallyExclusiveReactionsStrict(namespaceId string) [][]string {
	res, err := cs.GetMutuallyExclusiveReactions(namespaceId)
	if err != nil {
		log.Panicf("failed to get mutually exclusive reactions: %s", err)
	}
	return res
}

func (cs *ConfigurationStorage) GetMaxUniqueReactionsStrict(namespaceId string) int {
	res, err := cs.GetMaxUniqueReactions(namespaceId)
	if err != nil {
		log.Panicf("failed to get max unique reactions: %s", err)
	}
	return res
}

func (cs *ConfigurationStorage) ClearStrict(ctx context.Context) {
	err := cs.Clear(ctx)
	if err != nil {
		log.Panicf("failed to clear configuration storage: %s", err)
	}
}
