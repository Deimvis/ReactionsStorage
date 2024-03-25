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
	GetVisibleEntities() []Entity

	// Simulates switching to group of entities.
	// Calls Refresh for new visible entities.
	SwitchTopic(topicId string, userId string) Waitable

	// Returns true if the end of current topic isn't reached
	CanScroll() bool

	// Simulates scrolling.
	// Calls Refresh for new visible entities.
	// Returns an error if the end of current topic is reached.
	Scroll(userId string) (Waitable, error)

	// Simulates background reactions refresh.
	// Asynchronously updates reactions for visible entities.
	Refresh(userId string) Waitable

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
	app.topics = make([]Topic, len(topics))
	for i, t := range topics {
		app.topics[i] = t.Copy()
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

func (a *AppImpl) SwitchTopic(topicId string, userId string) Waitable {
	a.curTopicId = topicId
	a.curTopicPos = 0
	return a.Refresh(userId)
}

func (a *AppImpl) CanScroll() bool {
	return a.curTopicPos < len(a.getCurTopic().GetEntities())-a.visibleEntitiesCount
}

func (a *AppImpl) Scroll(userId string) (Waitable, error) {
	if !a.CanScroll() {
		return nil, fmt.Errorf("can't scroll: the end of current topic was reached")
	}
	a.curTopicPos = min(a.curTopicPos+a.visibleEntitiesCount, len(a.getCurTopic().GetEntities())-a.visibleEntitiesCount)
	return a.Refresh(userId), nil
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
			entityTsExpected := e.GetLastUpdateTs()
			resp, err := a.client.GetReactions(&req)
			if err != nil {
				a.logger.Warnw("Failed to get reactions", "err", err,
					"entity_id", e.GetId(), "user_id", userId)
				return
				// panic(fmt.Errorf("failed to get reactions: %w", err))
			}
			if resp.Code() == 200 {
				resp200, ok := resp.(*models.ReactionsGETResponse200)
				if !ok {
					panic(fmt.Errorf("failed to cast response to response200"))
				}
				reactionsCount := utils.Map(resp200.ReactionsCount, func(rc models.ReactionCount) ReactionCount { return ReactionCount(rc) })
				err := e.Update(reactionsCount, resp200.UserReactions.Reactions, WithLastUpdateTs(entityTsExpected))
				if err != nil {
					a.logger.Warnw("failed to refresh entity reactions since it was updated during request",
						"user_id", userId, "entity_id", e.GetId())
					return
				}
				// TODO: remove
				a.logger.Debugw("Refreshed reactions", "user_id", userId, "entity_id", e.GetId(), "my_reactions", e.GetMyReactionIds())
			} else {
				a.logger.Warnf("Get reactions status code: %d", resp.Code())
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
			a.logger.Warnw("Failed to add reaction", "err", err,
				"entity_id", e.GetId(), "reaction_id", reactionId, "user_id", userId)
			return
			// panic(fmt.Errorf("failed to add reaction: %w", err))
		}
		if resp.Code() == 200 {
			removedIds := GetConflictingReactionIds(e.GetNamespace(), reactionId)
			// TODO: remove
			a.logger.Debugw("Update my reactions", "user_id", userId, "entity_id", e.GetId(), "add", []string{reactionId}, "remove", removedIds)
			e.UpdateMyReactions([]string{reactionId}, removedIds)
		} else if resp.Code() == 403 {
			resp403, ok := resp.(*models.ReactionsPOSTResponse403)
			if !ok {
				panic("failed to cast response to response403")
			}
			a.logger.Warnw("Add reaction status code: 403 (constraint violated)", "error", resp403.Error)
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
		defer wg.Done()
		resp, err := a.client.RemoveReaction(&req)
		if err != nil {
			a.logger.Warnw("Failed to remove reaction", "err", err,
				"entity_id", e.GetId(), "reaction_id", reactionId, "user_id", userId)
			return
			// panic(fmt.Errorf("failed to remove reaction: %w", err))
		}
		if resp.Code() == 200 {
			e.RemoveMyReaction(reactionId)
		} else if resp.Code() == 403 {
			resp403, ok := resp.(*models.ReactionsDELETEResponse403)
			if !ok {
				panic("failed to cast response to response403")
			}
			a.logger.Warnw("Remove reaction status code: 403 (constraint violated)", "error", resp403.Error)
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
