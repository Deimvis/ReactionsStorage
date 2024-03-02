package models

type Configuration struct {
	Reactions    []Reaction    `filename:"reaction.yaml"`
	ReactionSets []ReactionSet `filename:"reaction_set.yaml"`
	// Namespaces   []Namespace
}

type Reaction struct {
	Id        string  `yaml:"id"`
	ShortName *string `yaml:"short_name"`
	Type      string  `yaml:"type"`
	Code      *string `yaml:"code"`
	Url       *string `yaml:"url"`
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
