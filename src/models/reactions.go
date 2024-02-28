package models

type UserReaction struct {
	UserId      string `json:"user_id" binding:"required"`
	ReactionId  string `json:"reaction_id" binding:"required"`
	NamespaceId string `json:"namespace_id" binding:"required"`
	EntityId    string `json:"entity_id" binding:"required"`
}

type ReactionCount struct {
	ReactionId string `json:"reaction_id"`
	Count      int    `json:"count"`
}

type UserReactionsWithinEntity struct {
	UserId    string   `json:"user_id"`
	Reactions []string `json:"reactions"`
}
