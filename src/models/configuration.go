package models

type Reaction struct {
	Id        string
	ShortName *string
	Type      string
	Code      *string
	Url       *string
}

type ReactionSet struct {
	Id          string
	ReactionIds []string
}

type Namespace struct {
	Id                         string
	ReactionSetId              string
	MaxUniqReactions           int
	MutuallyExclusiveReactions [][]string
}
