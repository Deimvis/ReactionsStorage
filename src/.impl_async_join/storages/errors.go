package storages

import "fmt"

type MaxUniqReactionsError struct{}

func (e *MaxUniqReactionsError) Error() string {
	return "max unique reactions reached for entity"
}

type ReactionAlreadyExistsError struct {
	userId     string
	reactionId string
}

func (e *ReactionAlreadyExistsError) Error() string {
	return fmt.Sprintf("reaction `%s` already exists for user `%s`", e.reactionId, e.userId)
}

type ConflictingReactionError struct {
	reactionId           string
	conflictingReactions []string
}

func (e *ConflictingReactionError) Error() string {
	return fmt.Sprintf("reaction `%s` conflicts with other reactions: %v", e.reactionId, e.conflictingReactions)
}
