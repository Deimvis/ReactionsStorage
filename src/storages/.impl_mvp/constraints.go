package storages

import (
	"context"
	"slices"
)

func checkAddUserReaction(ctx context.Context, userId string, reactionId string, uniqEntityReactions map[string]struct{}, uniqEntityUserReactions map[string]struct{}, maxUniqReactions int, mutExclReactions [][]string) error {
	_, alreadyExists := uniqEntityUserReactions[reactionId]
	if alreadyExists {
		return &ReactionAlreadyExistsError{userId, reactionId}
	}

	uniqReactions := len(uniqEntityReactions)
	_, isNewEntityReaction := uniqEntityReactions[reactionId]
	if isNewEntityReaction {
		uniqReactions += 1
	}
	if uniqReactions >= maxUniqReactions {
		return &MaxUniqReactionsError{}
	}

	conflictingReactions := getConflictingReactionIds(reactionId, uniqEntityUserReactions, mutExclReactions)
	if len(conflictingReactions) > 0 {
		return &ConflictingReactionError{reactionId, conflictingReactions}
	}

	return nil
}

func getConflictingReactionIds(reactionId string, uniqEntityUserReactions map[string]struct{}, mutExclReactions [][]string) []string {
	var conflictingReactions []string
	for _, conflictingGroup := range mutExclReactions {
		ind := slices.Index(conflictingGroup, reactionId)
		if ind == -1 {
			continue
		}
		for _, other := range conflictingGroup {
			_, has := uniqEntityUserReactions[other]
			if other != reactionId && has {
				conflictingReactions = append(conflictingReactions, other)
			}
		}
	}
	return conflictingReactions
}
