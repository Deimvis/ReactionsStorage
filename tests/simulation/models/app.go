package models

import (
	"fmt"
	"math/rand"
	"slices"

	"github.com/Deimvis/reactionsstorage/src/models"
	rs "github.com/Deimvis/reactionsstorage/tests/simulation/rs_client"
	"github.com/Deimvis/reactionsstorage/tests/simulation/utils"
)

// App encapsulates application functionality
type App interface {
	GetCurrentTopicId() string
	GetAvailableTopicIds() []string
	SwitchTopic(topicId string)
	CanScroll() bool // returns true if the end of current topic isn't reached
	Scroll() error   // returns an error if the end of current topic is reached
	GetVisibleEntities() []Entity

	// Simulates click on reaction.
	// Asynchronously sends request and updates given entity.
	AddReaction(e Entity, userId string, reactionId string)

	// Simulates click on existing reaction.
	// Asynchronously sends request and updates given entity.
	RemoveReaction(e Entity, userId string, reactionId string)
}

type AppImpl struct {
	client               rs.Client
	topics               []Topic
	visibleEntitiesCount int
	curTopicId           string
	curTopicPos          int
}

func NewApp(client rs.Client, topics []Topic, visibleEntitiesCount int) App {
	if len(topics) == 0 {
		panic(fmt.Errorf("no topics"))
	}
	var app AppImpl
	app.client = client
	app.visibleEntitiesCount = visibleEntitiesCount
	copy(app.topics, topics)
	for i := range app.topics {
		app.topics[i].ShuffleEntities()
	}
	app.curTopicId = app.topics[rand.Intn(len(app.topics))].GetId()
	return &app
}

func (a *AppImpl) GetCurrentTopicId() string {
	return a.curTopicId
}

func (a *AppImpl) GetAvailableTopicIds() []string {
	return utils.Map(a.topics, func(t Topic) string { return t.GetId() })
}

func (a *AppImpl) SwitchTopic(topicId string) {
	a.curTopicId = topicId
	a.curTopicPos = 0
}

func (a *AppImpl) CanScroll() bool {
	return a.curTopicPos < len(a.getCurTopic().GetEntities())-a.visibleEntitiesCount
}

func (a *AppImpl) Scroll() error {
	if !a.CanScroll() {
		return fmt.Errorf("can't scroll: the end of current topic was reached")
	}
	a.curTopicPos = min(a.curTopicPos+a.visibleEntitiesCount, len(a.getCurTopic().GetEntities())-a.visibleEntitiesCount)
	return nil
}

func (a *AppImpl) GetVisibleEntities() []Entity {
	return a.getCurTopic().GetEntities()[a.curTopicPos : a.curTopicPos+a.visibleEntitiesCount]
}

func (a *AppImpl) AddReaction(e Entity, userId string, reactionId string) {
	utils.AssertPtr(e)

	var req models.ReactionsPOSTRequest
	// force -- automatically remove conflicting reactions
	force := true
	req.Query.Force = &force
	req.Body.NamespaceId = e.GetNamespace().GetId()
	req.Body.EntityId = e.GetId()
	req.Body.UserId = userId
	req.Body.ReactionId = reactionId
	go func() {
		resp, err := a.client.AddReaction(&req)
		if err != nil {
			panic(fmt.Errorf("failed to add reaction: %w", err))
		}
		if resp.Code() == 200 {
			removedIds := GetConflictingReactionIds(e.GetNamespace(), reactionId)
			e.UpdateMyReactions([]string{reactionId}, removedIds)
		}
	}()
}

func (a *AppImpl) RemoveReaction(e Entity, userId string, reactionId string) {
	utils.AssertPtr(e)

	var req models.ReactionsDELETERequest
	req.Body.NamespaceId = e.GetNamespace().GetId()
	req.Body.EntityId = e.GetId()
	req.Body.UserId = userId
	req.Body.ReactionId = reactionId
	go func() {
		resp, err := a.client.RemoveReaction(&req)
		if err != nil {
			panic(fmt.Errorf("failed to remove reaction: %w", err))
		}
		if resp.Code() == 200 {
			e.RemoveMyReaction(reactionId)
		}
	}()
}

func (a *AppImpl) getCurTopic() Topic {
	return a.getTopic(a.curTopicId)
}

func (a *AppImpl) getTopic(topicId string) Topic {
	return a.topics[slices.IndexFunc(a.topics, func(t Topic) bool { return t.GetId() == topicId })]
}
