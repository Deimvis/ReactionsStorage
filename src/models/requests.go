package models

import "net/http"

// type ReactionsGETRequest =

type Request interface {
	Method() string
	Path() string
	QueryString() string
	Header() http.Header
	BodyRaw() []byte
}

type ReactionsGETRequest struct {
	Query struct {
		NamespaceId string `query:"namespace_id"`
		EntityId    string `query:"entity_id"`
		UserId      string `query:"user_id"`
	}
}

type ReactionsPOSTRequest struct {
	Query struct {
		Force *bool `query:"force"`
	}
	Body UserReaction
}

type ReactionsDELETERequest struct {
	Body UserReaction
}

type ConfigurationPOSTRequest struct {
	Headers http.Header
	Body    Configuration
}

type NamespaceGETRequest struct {
	Query struct {
		NamespaceId string `query:"namespace_id"`
	}
}

type AvailableReactionsGETRequest struct {
	Query struct {
		NamespaceId string `query:"namespace_id"`
	}
}
