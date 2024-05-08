package storages

import (
	"context"
	"slices"

	"github.com/Deimvis/reactionsstorage/src/utils"
)

func checkAddUserReaction(ctx context.Context, userId string, reactionId string, reactionsCount map[string]int, uniqUserReactions []string, maxUniqReactions int, mutExclReactions [][]string) error {
	if utils.Contains(uniqUserReactions, reactionId) {
		return &ReactionAlreadyExistsError{userId, reactionId}
	}

	// NOTE: conflictign reactions should be check first.
	// In case with `force` flag, removing conflicting reactions can
	// decrease number of unique reactions and allow to pass max_uniq_reactions constraint.
	conflictingReactions := getConflictingReactionIds(reactionId, uniqUserReactions, mutExclReactions)
	if len(conflictingReactions) > 0 {
		return &ConflictingReactionError{reactionId, conflictingReactions}
	}

	uniqReactions := len(reactionsCount)
	_, isExistingEntityReaction := reactionsCount[reactionId]
	if !isExistingEntityReaction {
		uniqReactions += 1
	}
	if uniqReactions > maxUniqReactions {
		return &MaxUniqReactionsError{}
	}

	return nil
}

func getConflictingReactionIds(reactionId string, uniqUserReactions []string, mutExclReactions [][]string) []string {
	var conflictingReactions []string
	for _, conflictingGroup := range mutExclReactions {
		ind := slices.Index(conflictingGroup, reactionId)
		if ind == -1 {
			continue
		}
		for _, other := range conflictingGroup {
			if other != reactionId && utils.Contains(uniqUserReactions, other) {
				conflictingReactions = append(conflictingReactions, other)
			}
		}
	}
	return conflictingReactions
}
