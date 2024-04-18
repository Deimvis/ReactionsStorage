package models

import (
	"errors"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/tests/simulation/src/utils"
)

type ReactionCount = models.ReactionCount

// thread-safe
type Entity interface {
	GetId() string
	GetNamespace() Namespace
	GetReactionsCount() []ReactionCount
	GetMyReactionIds() []string
	GetLastUpdateTs() int64

	// Atomically update state
	Update(rcs []ReactionCount, myReactions []string, conditions ...Condition) error
	// Adds reactions atomically.
	// Returns false if reaction was already added.
	AddMyReaction(reactionId string) bool
	// Removes reactions atomically.
	// Returns false if reaction was already removed.
	RemoveMyReaction(reactionId string) bool
	// First adds and then removes reactions atomically.
	// Returns statuses (see AddMyReaction and RemoveMyReaction).
	UpdateMyReactions(addIds []string, removeIds []string, conditions ...Condition) ([]bool, error)

	getLastUpdateTs() int64
}

func NewEntity(id string, namespace Namespace) Entity {
	return &EntityImpl{id: id, namespace: namespace, lastUpdateTs: time.Now().Unix()}
}

type Condition func(e Entity) bool

func WithLastUpdateTs(ts int64) Condition {
	return func(e Entity) bool {
		return e.getLastUpdateTs() == ts
	}
}

// thread safe
type EntityImpl struct {
	id             string
	namespace      Namespace
	reactionsCount []ReactionCount
	myReactionIds  []string

	lastUpdateTs int64
	m            sync.Mutex
}

func (e *EntityImpl) GetId() string {
	e.m.Lock()
	defer e.m.Unlock()
	return e.id
}

func (e *EntityImpl) GetNamespace() Namespace {
	e.m.Lock()
	defer e.m.Unlock()
	return e.namespace
}

func (e *EntityImpl) GetReactionsCount() []ReactionCount {
	e.m.Lock()
	defer e.m.Unlock()
	return e.reactionsCount
}

func (e *EntityImpl) GetMyReactionIds() []string {
	e.m.Lock()
	defer e.m.Unlock()
	return e.myReactionIds
}

func (e *EntityImpl) GetLastUpdateTs() int64 {
	e.m.Lock()
	defer e.m.Unlock()
	return e.getLastUpdateTs()
}

func (e *EntityImpl) Update(rcs []ReactionCount, myReactions []string, conditions ...Condition) error {
	e.m.Lock()
	defer e.m.Unlock()
	if !e.fulfills(conditions) {
		return errors.New("some conditions weren't fulfilled")
	}
	defer e.updateTs()
	e.reactionsCount = rcs
	e.myReactionIds = myReactions
	return nil
}

func (e *EntityImpl) AddMyReaction(reactionId string) bool {
	res, err := e.UpdateMyReactions([]string{reactionId}, nil)
	if err != nil {
		panic("no error was expected")
	}
	if len(res) != 1 {
		panic("bug: incorrect length of res")
	}
	return res[0]
}

func (e *EntityImpl) RemoveMyReaction(reactionId string) bool {
	res, err := e.UpdateMyReactions(nil, []string{reactionId})
	if err != nil {
		panic("no error was expected")
	}
	if len(res) != 1 {
		panic("bug: incorrect length of res")
	}
	return res[0]
}

func (e *EntityImpl) UpdateMyReactions(addIds []string, removeIds []string, conditions ...Condition) ([]bool, error) {
	e.m.Lock()
	defer e.m.Unlock()
	if !e.fulfills(conditions) {
		return nil, errors.New("some conditions weren't fulfilled")
	}
	defer e.updateTs()
	var res []bool
	for _, id := range addIds {
		res = append(res, e.addMyReactionUnsafe(id))
	}
	for _, id := range removeIds {
		res = append(res, e.removeMyReactionUnsafe(id))
	}
	return res, nil
}

func (e *EntityImpl) fulfills(conditions []Condition) bool {
	for _, cond := range conditions {
		if !cond(e) {
			return false
		}
	}
	return true
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

func (e *EntityImpl) getLastUpdateTs() int64 {
	return e.lastUpdateTs
}

func (e *EntityImpl) updateTs() {
	e.lastUpdateTs = time.Now().Unix()
}
