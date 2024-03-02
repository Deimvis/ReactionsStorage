package models

// type ReactionsGETRequest =

type Request interface {
	Method() string
	Path() string
	QueryString() string
	BodyRaw() []byte
}

type ReactionsGETRequest struct {
	Query struct {
		NamespaceId string `query:"namespace_id"`
		EntityId    string
		UserId      string
	}
}

type ReactionsPOSTRequest struct {
	Query struct {
		Force *bool
	}
	Body  UserReaction
}

type ReactionsDELETERequest struct {
	Body UserReaction
}

type ConfiguratinPOSTRequest struct {
	Body Configuration
}
