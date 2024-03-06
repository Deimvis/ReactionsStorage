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

type ReactionsDELETEResponse200 struct {
	Status string `json:"status"`
}

type ConfigurationPOSTResponse200 struct {
	Status string `json:"status"`
}

type AvailableReactionsGETResponse200 struct {
	Reactions []Reaction `json:"reactions"`
}

type NamespaceGETResponse200 struct {
	Namespace Namespace `json:"namespace"`
}

type ReactionsPOSTResponse403 ErrorResponse
type ReactionsDELETEResponse403 ErrorResponse
type ConfigurationPOSTResponse415 ErrorResponse
type ConfigurationPOSTResponse422 ErrorResponse
type NamespaceGETResponse404 ErrorResponse
type AvailableReactionsGETResponse404 ErrorResponse

type ErrorResponse struct {
	Error string `json:"error"`
}
