package fake

import (
	"github.com/Deimvis/reactionsstorage/src/models"
)

const (
	MaxUniqReactions = 5
	NamespaceId      = "namespace"
	ReactionId       = ReactionId1
	ReactionId1      = "reaction1"
	ReactionId2      = "reaction2"
	ReactionId3      = "reaction3"
	FreeReactionId1  = "free_reaction1" // no constraints applied
	FreeReactionId2  = "free_reaction2" // no constraints applied
	FreeReactionId3  = "free_reaction3" // no constraints applied
	FreeReactionId4  = "free_reaction4" // no constraints applied
	FreeReactionId5  = "free_reaction5" // no constraints applied
	FreeReactionId6  = "free_reaction6" // no constraints applied
	ReactionSetId    = "reaction_set"
	ReactionType     = "unicode"
	EntityId         = EntityId1
	EntityId1        = "entity1"
	EntityId2        = "entity2"
	EntityId3        = "entity3"
	UserId           = UserId1
	UserId1          = "user1"
	UserId2          = "user2"
	UserId3          = "user3"
)

var (
	ReactionCode               = "😃"
	MutuallyExclusiveReactions = [][]string{{ReactionId1, ReactionId2}, {ReactionId2, ReactionId3}}
)

var (
	Reaction  = Reaction1
	Reaction1 = models.Reaction{
		Id:   ReactionId1,
		Type: ReactionType,
		Code: &ReactionCode,
	}
	Reaction2 = models.Reaction{
		Id:   ReactionId2,
		Type: ReactionType,
		Code: &ReactionCode,
	}
	Reaction3 = models.Reaction{
		Id:   ReactionId3,
		Type: ReactionType,
		Code: &ReactionCode,
	}

	FreeReaction  = FreeReaction1
	FreeReaction1 = models.Reaction{
		Id:   FreeReactionId1,
		Type: ReactionType,
		Code: &ReactionCode,
	}
	FreeReaction2 = models.Reaction{
		Id:   FreeReactionId2,
		Type: ReactionType,
		Code: &ReactionCode,
	}
	FreeReaction3 = models.Reaction{
		Id:   FreeReactionId3,
		Type: ReactionType,
		Code: &ReactionCode,
	}
	FreeReaction4 = models.Reaction{
		Id:   FreeReactionId4,
		Type: ReactionType,
		Code: &ReactionCode,
	}
	FreeReaction5 = models.Reaction{
		Id:   FreeReactionId5,
		Type: ReactionType,
		Code: &ReactionCode,
	}
	FreeReaction6 = models.Reaction{
		Id:   FreeReactionId6,
		Type: ReactionType,
		Code: &ReactionCode,
	}
)

var ReactionSet = models.ReactionSet{
	Id: ReactionSetId,
	ReactionIds: []string{
		ReactionId1, ReactionId2, ReactionId3,
		FreeReactionId1, FreeReactionId2, FreeReactionId3,
		FreeReactionId4, FreeReactionId5, FreeReactionId6,
	},
}

var Namespace = models.Namespace{
	Id:                         NamespaceId,
	ReactionSetId:              ReactionSetId,
	MaxUniqReactions:           MaxUniqReactions,
	MutuallyExclusiveReactions: MutuallyExclusiveReactions,
}
