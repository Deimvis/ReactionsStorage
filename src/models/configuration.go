package models

type Configuration struct {
	Reactions    []Reaction    `json:"reactions" filename:"reaction.yaml"`
	ReactionSets []ReactionSet `json:"reaction_sets" filename:"reaction_set.yaml"`
	Namespaces   []Namespace   `json:"namespaces" filename:"namespace.yaml"`
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
