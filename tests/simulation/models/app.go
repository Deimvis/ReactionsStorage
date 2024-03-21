package models

import (
	"fmt"
	"math/rand"
	"slices"
	"sync"

	"go.uber.org/zap"

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
	Refresh(userId string) Waitable // asynchronously updates reactions for visible entities

	// Simulates click on reaction.
	// Asynchronously sends request and updates given entity.
	AddReaction(e Entity, userId string, reactionId string) Waitable

	// Simulates click on existing reaction.
	// Asynchronously sends request and updates given entity.
	RemoveReaction(e Entity, userId string, reactionId string) Waitable
}

type AppImpl struct {
	client               rs.Client
	topics               []Topic
	visibleEntitiesCount int
	curTopicId           string
	curTopicPos          int

	logger *zap.SugaredLogger
}

func NewApp(client rs.Client, topics []Topic, visibleEntitiesCount int, logger *zap.SugaredLogger) App {
	if len(topics) == 0 {
		panic(fmt.Errorf("no topics"))
	}
	var app AppImpl
	app.client = client
	app.visibleEntitiesCount = visibleEntitiesCount
	app.topics = append(app.topics, topics...)
	for i := range app.topics {
		app.topics[i].ShuffleEntities()
	}
	app.curTopicId = app.topics[rand.Intn(len(app.topics))].GetId()
	app.logger = logger
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

func (a *AppImpl) Refresh(userId string) Waitable {
	wg := sync.WaitGroup{}
	for _, e := range a.GetVisibleEntities() {
		utils.AssertPtr(e)

		var req models.ReactionsGETRequest
		req.Query.NamespaceId = e.GetNamespace().GetId()
		req.Query.EntityId = e.GetId()
		req.Query.UserId = userId

		wg.Add(1)
		go func(e Entity) {
			defer wg.Done()
			resp, err := a.client.GetReactions(&req)
			if err != nil {
				panic(fmt.Errorf("failed to get reactions: %w", err))
			}
			if resp.Code() == 200 {
				resp200, ok := resp.(*models.ReactionsGETResponse200)
				if !ok {
					panic(fmt.Errorf("failed to cast response to response200"))
				}
				reactionsCount := utils.Map(resp200.ReactionsCount, func(rc models.ReactionCount) ReactionCount { return ReactionCount(rc) })
				e.Update(reactionsCount, resp200.UserReactions.Reactions)
			}
		}(e)
	}
	return &wg
}

func (a *AppImpl) AddReaction(e Entity, userId string, reactionId string) Waitable {
	utils.AssertPtr(e)
	wg := sync.WaitGroup{}

	var req models.ReactionsPOSTRequest
	// force -- automatically remove conflicting reactions
	force := true
	req.Query.Force = &force
	req.Body.NamespaceId = e.GetNamespace().GetId()
	req.Body.EntityId = e.GetId()
	req.Body.UserId = userId
	req.Body.ReactionId = reactionId

	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := a.client.AddReaction(&req)
		if err != nil {
			panic(fmt.Errorf("failed to add reaction: %w", err))
		}
		if resp.Code() == 200 {
			removedIds := GetConflictingReactionIds(e.GetNamespace(), reactionId)
			e.UpdateMyReactions([]string{reactionId}, removedIds)
		} else {
			a.logger.Warnf("Add reaction status code: %d", resp.Code())
		}
	}()
	return &wg
}

func (a *AppImpl) RemoveReaction(e Entity, userId string, reactionId string) Waitable {
	utils.AssertPtr(e)
	wg := sync.WaitGroup{}

	var req models.ReactionsDELETERequest
	req.Body.NamespaceId = e.GetNamespace().GetId()
	req.Body.EntityId = e.GetId()
	req.Body.UserId = userId
	req.Body.ReactionId = reactionId

	wg.Add(1)
	go func() {
		resp, err := a.client.RemoveReaction(&req)
		if err != nil {
			panic(fmt.Errorf("failed to remove reaction: %w", err))
		}
		if resp.Code() == 200 {
			e.RemoveMyReaction(reactionId)
		} else {
			a.logger.Warnf("Remove reaction status code: %d", resp.Code())
		}
	}()
	return &wg
}

func (a *AppImpl) getCurTopic() Topic {
	return a.getTopic(a.curTopicId)
}

func (a *AppImpl) getTopic(topicId string) Topic {
	return a.topics[slices.IndexFunc(a.topics, func(t Topic) bool { return t.GetId() == topicId })]
}

type Waitable interface {
	Wait()
}
