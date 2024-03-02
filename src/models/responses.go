package models

type Response interface {
	Code() int
}

type ReactionsGETResponse200 struct {
	EntityId       string                    `json:"entity_id"`
	ReactionsCount []ReactionCount           `json:"reactions_count"`
	UserReactions  UserReactionsWithinEntity `json:"user_reactions"`
}

type ReactionsPOSTResponse200 struct {
	Status string `json:"status"`
}

type ReactionsPOSTResponse403 struct {
	Error string `json:"error"`
}

type ReactionsDELETEResponse200 struct {
	Status string `json:"status"`
}

type ReactionsDELETEResponse403 struct {
	Error string `json:"error"`
}

type ConfigurationPOSTResponse200 struct {
	Status string `json:"status"`
}

type ConfigurationPOSTResponse422 struct {
	Error string `json:"error"`
}
