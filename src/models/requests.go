package models

// type ReactionsGETRequest =

type ReactionsGETRequest struct {
	Query struct {
		NamespaceId string
		EntityId    string
		UserId      string
	}
}

type ReactionsPOSTRequest struct {
	Body  UserReaction
	Query struct {
		Force bool
	}
}

type ReactionsDELETERequest struct {
	Body UserReaction
}
