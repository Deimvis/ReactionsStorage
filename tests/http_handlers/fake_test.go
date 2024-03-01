package http_handlers_test

import (
	"github.com/Deimvis/reactionsstorage/src/models"
)

const (
	fakeEntityId         = "entity"
	fakeMaxUniqReactions = 10
	fakeNamespaceId      = "namespace"
	fakeReactionId       = "reaction"
	fakeReactionId2      = "reaction2"
	fakeReactionId3      = "reaction3"
	fakeReactionSetId    = "reaction_set"
	fakeReactionType     = "unicode"
	fakeUserId           = "user"
)

var (
	fakeReactionCode               = "😃"
	fakeMutuallyExclusiveReactions = [][]string{{fakeReactionId, fakeReactionId2}, {fakeReactionId2, fakeReactionId3}}
)

var fakeReaction = models.Reaction{
	Id:   fakeReactionId,
	Type: fakeReactionType,
	Code: &fakeReactionCode,
}

var fakeReaction2 = models.Reaction{
	Id:   fakeReactionId2,
	Type: fakeReactionType,
	Code: &fakeReactionCode,
}

var fakeReaction3 = models.Reaction{
	Id:   fakeReactionId3,
	Type: fakeReactionType,
	Code: &fakeReactionCode,
}

var fakeReactionSet = models.ReactionSet{
	Id:          fakeReactionSetId,
	ReactionIds: []string{fakeReactionId, fakeReactionId2, fakeReactionId3},
}

var fakeNamespace = models.Namespace{
	Id:                         fakeNamespaceId,
	ReactionSetId:              fakeReactionSetId,
	MaxUniqReactions:           fakeMaxUniqReactions,
	MutuallyExclusiveReactions: fakeMutuallyExclusiveReactions,
}
