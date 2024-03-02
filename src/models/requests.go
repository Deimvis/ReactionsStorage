package models

// type ReactionsGETRequest =

type Request interface {
	Method() string
	Path() string
	QueryString() string
	BodyJSON() []byte
}

type ReactionsGETRequest struct {
	Query struct {
		NamespaceId string `query:"namespace_id"`
		EntityId    string
		UserId      string
	}
}

type ReactionsPOSTRequest struct {
	Body  UserReaction
	Query struct {
		Force *bool
	}
}

type ReactionsDELETERequest struct {
	Body UserReaction
}
