package models

import (
	"fmt"
	"slices"
)

type Configuration struct {
	Reactions    []Reaction    `yaml:"reaction" json:"reaction" filename:"reaction.yaml"`
	ReactionSets []ReactionSet `yaml:"reaction_set" json:"reaction_set" filename:"reaction_set.yaml"`
	Namespaces   []Namespace   `yaml:"namespace" json:"namespace" filename:"namespace.yaml"`
}

type Reaction struct {
	Id        string  `yaml:"id" db:"id"`
	ShortName *string `yaml:"short_name" db:"short_name"`
	Type      string  `yaml:"type" db:"type"`
	Code      *string `yaml:"code" db:"code"`
	Url       *string `yaml:"url" db:"url"`
}

type ReactionSet struct {
	Id          string   `yaml:"id" json:"id" db:"id"`
	ReactionIds []string `yaml:"reaction_ids" json:"reaction_ids" db:"reaction_ids"`
}

type Namespace struct {
	Id                         string     `yaml:"id" json:"id" db:"id"`
	ReactionSetId              string     `yaml:"reaction_set_id" json:"reaction_set_id" db:"reaction_set_id"`
	MaxUniqReactions           int        `yaml:"max_uniq_reactions" json:"max_uniq_reactions" db:"max_uniq_reactions"`
	MutuallyExclusiveReactions [][]string `yaml:"mutually_exclusive_reactions" json:"mutually_exclusive_reactions" db:"mutually_exclusive_reactions"`
}

func CheckCorrectness(c *Configuration) error {
	{
		idsSeen := make(map[string]struct{})
		for _, r := range c.Reactions {
			switch r.Type {
			case "unicode":
				if r.Code == nil || len(*r.Code) == 0 {
					return fmt.Errorf("reaction(id=`%s`) has unicode type, but no `code` specified", r.Id)
				}
			case "custom":
				if r.Url == nil || len(*r.Url) == 0 {
					return fmt.Errorf("reaction(id=`%s`) has custom type, but no `url` specified", r.Id)
				}
			default:
				return fmt.Errorf("reaction(id=`%s`) has incorrect type: %s", r.Id, r.Type)
			}
			if _, ok := idsSeen[r.Id]; ok {
				return fmt.Errorf("found more than 1 reaction with id=`%s`", r.Id)
			}
			idsSeen[r.Id] = struct{}{}
		}
	}

	{
		idsSeen := make(map[string]struct{})
		for _, rs := range c.ReactionSets {
			for _, rId := range rs.ReactionIds {
				if slices.IndexFunc(c.Reactions, func(r Reaction) bool { return r.Id == rId }) == -1 {
					return fmt.Errorf("reaction_set(id=`%s`) has reaction(id=`%s`) that doesn't exist", rs.Id, rId)
				}
			}
			if _, ok := idsSeen[rs.Id]; ok {
				return fmt.Errorf("found more than 1 reaction set with id=`%s`", rs.Id)
			}
			idsSeen[rs.Id] = struct{}{}
		}
	}

	{
		idsSeen := make(map[string]struct{})
		for _, n := range c.Namespaces {
			if slices.IndexFunc(c.ReactionSets, func(rs ReactionSet) bool { return rs.Id == n.ReactionSetId }) == -1 {
				return fmt.Errorf("namespace(id=`%s`) has reaction_set(id=`%s`) that doesn't exist", n.Id, n.ReactionSetId)
			}
			if n.MaxUniqReactions < 0 {
				return fmt.Errorf("namespace(id=`%s`) has negative max_uniq_reactions: %d", n.Id, n.MaxUniqReactions)
			}
			for _, conflictingGroup := range n.MutuallyExclusiveReactions {
				for _, rId := range conflictingGroup {
					if slices.IndexFunc(c.Reactions, func(r Reaction) bool { return r.Id == rId }) == -1 {
						return fmt.Errorf("namespace(id=`%s`) has mutually exclusive reaction group with reaction(id=`%s`) that doesn't exist (group=`%v`)", n.Id, rId, conflictingGroup)
					}
				}
			}
			if _, ok := idsSeen[n.Id]; ok {
				return fmt.Errorf("found more than 1 namespace set with id=`%s`", n.Id)
			}
			idsSeen[n.Id] = struct{}{}
		}
	}
	return nil
}
