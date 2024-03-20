package fake

import (
	"github.com/Deimvis/reactionsstorage/src/models"
)

const (
	EntityId         = "entity"
	MaxUniqReactions = 10
	NamespaceId      = "namespace"
	ReactionId       = ReactionId1
	ReactionId1      = "reaction1"
	ReactionId2      = "reaction2"
	ReactionId3      = "reaction3"
	ReactionSetId    = "reaction_set"
	ReactionType     = "unicode"
	UserId           = "user"
)

var (
	ReactionCode               = "😃"
	MutuallyExclusiveReactions = [][]string{{ReactionId1, ReactionId2}, {ReactionId2, ReactionId3}}
)

var Reaction = models.Reaction{
	Id:   ReactionId,
	Type: ReactionType,
	Code: &ReactionCode,
}

var Reaction2 = models.Reaction{
	Id:   ReactionId2,
	Type: ReactionType,
	Code: &ReactionCode,
}

var Reaction3 = models.Reaction{
	Id:   ReactionId3,
	Type: ReactionType,
	Code: &ReactionCode,
}

var ReactionSet = models.ReactionSet{
	Id:          ReactionSetId,
	ReactionIds: []string{ReactionId1, ReactionId2, ReactionId3},
}

var Namespace = models.Namespace{
	Id:                         NamespaceId,
	ReactionSetId:              ReactionSetId,
	MaxUniqReactions:           MaxUniqReactions,
	MutuallyExclusiveReactions: MutuallyExclusiveReactions,
}
