package models

import (
	"sync"

	"go.uber.org/zap"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/tests/simulation/utils"
)

type ReactionCount = models.ReactionCount

// thread-safe
type Entity interface {
	GetId() string
	GetNamespace() Namespace
	GetReactionsCount() []ReactionCount
	GetMyReactionIds() []string

	// Atomically update state
	Update(rcs []ReactionCount, myReactions []string)
	// Adds reactions atomically.
	// Returns false if reaction was already added.
	AddMyReaction(reactionId string) bool
	// Removes reactions atomically.
	// Returns false if reaction was already removed.
	RemoveMyReaction(reactionId string) bool
	// First adds and then removes reactions atomically.
	// Returns statuses (see AddMyReaction and RemoveMyReaction).
	UpdateMyReactions(addIds []string, removeIds []string) []bool
}

func NewEntity(id string, namespace Namespace) Entity {
	return &EntityImpl{id: id, namespace: namespace}
}

// thread safe
type EntityImpl struct {
	id             string
	namespace      Namespace
	reactionsCount []ReactionCount
	myReactionIds  []string

	m sync.Mutex
}

func (e *EntityImpl) GetId() string {
	return e.id
}

func (e *EntityImpl) GetNamespace() Namespace {
	return e.namespace
}

func (e *EntityImpl) GetReactionsCount() []ReactionCount {
	return e.reactionsCount
}

func (e *EntityImpl) GetMyReactionIds() []string {
	return e.myReactionIds
}

func (e *EntityImpl) Update(rcs []ReactionCount, myReactions []string) {
	e.m.Lock()
	defer e.m.Unlock()
	e.reactionsCount = rcs
	e.myReactionIds = myReactions
}

func (e *EntityImpl) AddMyReaction(reactionId string) bool {
	e.m.Lock()
	defer e.m.Unlock()
	return e.addMyReactionUnsafe(reactionId)
}

func (e *EntityImpl) RemoveMyReaction(reactionId string) bool {
	e.m.Lock()
	defer e.m.Unlock()
	return e.removeMyReactionUnsafe(reactionId)
}

func (e *EntityImpl) UpdateMyReactions(addIds []string, removeIds []string) []bool {
	e.m.Lock()
	defer e.m.Unlock()
	var res []bool
	for _, id := range addIds {
		res = append(res, e.addMyReactionUnsafe(id))
	}
	for _, id := range removeIds {
		res = append(res, e.removeMyReactionUnsafe(id))
	}
	return res
}

func (e *EntityImpl) addMyReactionUnsafe(reactionId string) bool {
	if utils.Contains(e.myReactionIds, reactionId) {
		return false
	}

	e.myReactionIds = append(e.myReactionIds, reactionId)
	isNew := true
	for i := range e.reactionsCount {
		if e.reactionsCount[i].ReactionId == reactionId {
			isNew = false
			e.reactionsCount[i].Count++
		}
	}
	if isNew {
		rc := ReactionCount{ReactionId: reactionId, Count: 1}
		e.reactionsCount = append(e.reactionsCount, rc)
	}
	return true
}

func (e *EntityImpl) removeMyReactionUnsafe(reactionId string) bool {
	if !utils.Contains(e.myReactionIds, reactionId) {
		return false
	}

	utils.FilterIn(&e.myReactionIds, func(id string) bool { return id == reactionId })
	removed := false
	for i := range e.reactionsCount {
		if e.reactionsCount[i].ReactionId == reactionId && e.reactionsCount[i].Count >= 1 {
			removed = true
			e.reactionsCount[i].Count--
		}
	}
	if !removed {
		zap.L().Warn("inconsistency happened: reaction was found in myReactions, but wasn't present in reactionsCount")
	}
	return true
}
