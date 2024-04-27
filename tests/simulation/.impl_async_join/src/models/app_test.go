package models

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/tests/fake"
	"github.com/Deimvis/reactionsstorage/tests/setup"
	"github.com/Deimvis/reactionsstorage/tests/simulation/src/utils"
)

func TestApp_Refresh(t *testing.T) {
	n := NewNamespace(fake.NamespaceId, setup.RSClient)
	app := NewApp(setup.RSClient, []Topic{NewTopic("", n, 1, true)}, 1, setup.Logger)
	e := app.GetVisibleEntities()[0]

	testCases := []struct {
		initial state
		userId  string
		expect  state
	}{
		{
			state{nil},
			fake.UserId,
			state{nil},
		},
		{

			state{nil},
			fake.UserId,
			StateBuilder().ForEntity(e).UserRs(fake.UserId, []string{fake.ReactionId}).Build(),
		},
		{
			state{nil},
			fake.UserId1,
			StateBuilder().ForEntity(e).UserRs(fake.UserId2, []string{fake.ReactionId}).Build(),
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			utils.AssertPtr(e)
			tc.initial.ApplyToEntity(e, tc.userId)
			tc.expect.ApplyToStorage()
			app.Refresh(tc.userId).Wait()
			tc.expect.RequireIsActualForEntity(t, e, tc.userId)
		})
	}
	setup.ClearUserReactions()
}

func TestApp_AddReaction(t *testing.T) {
	n := NewNamespace(fake.NamespaceId, setup.RSClient)
	app := NewApp(setup.RSClient, []Topic{NewTopic("", n, 1, true)}, 1, setup.Logger)
	e := app.GetVisibleEntities()[0]

	testCases := []struct {
		initial    state
		userId     string
		reactionId string
		expect     state
	}{
		{
			state{nil},
			fake.UserId,
			fake.ReactionId,
			StateBuilder().ForEntity(e).UserRs(fake.UserId, []string{fake.ReactionId}).Build(),
		},
		{
			StateBuilder().ForEntity(e).UserRs(fake.UserId1, []string{fake.ReactionId1}).Build(),
			fake.UserId1,
			fake.ReactionId3,
			StateBuilder().ForEntity(e).UserRs(fake.UserId1, []string{fake.ReactionId1, fake.ReactionId3}).Build(),
		},
		{
			StateBuilder().ForEntity(e).UserRs(fake.UserId1, []string{fake.ReactionId}).Build(),
			fake.UserId2,
			fake.ReactionId,
			StateBuilder().ForEntity(e).UserRs(fake.UserId1, []string{fake.ReactionId}).UserRs(fake.UserId2, []string{fake.ReactionId}).Build(),
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			utils.AssertPtr(e)
			tc.initial.Apply(e, tc.userId)
			app.AddReaction(e, tc.userId, tc.reactionId).Wait()
			setup.RS.RefreshEntityReactions(context.Background())
			tc.expect.RequireIsActual(t, e, tc.userId)
		})
	}
}

func StateBuilder() *stateBuilder {
	sb := &stateBuilder{}
	sb.entityId2entity = make(map[string]Entity)
	sb.entityId2usersReactions = make(map[string][]models.UserReactionsWithinEntity)
	return sb
}

type stateBuilder struct {
	curEntityId             string
	entityId2entity         map[string]Entity
	entityId2usersReactions map[string][]models.UserReactionsWithinEntity
}

func (sb *stateBuilder) ForEntity(entity Entity) *stateBuilder {
	sb.curEntityId = entity.GetId()
	sb.entityId2entity[entity.GetId()] = entity
	return sb
}

func (sb *stateBuilder) UserRs(userId string, reactionIds []string) *stateBuilder {
	urs := sb.entityId2usersReactions
	urs[sb.curEntityId] = append(urs[sb.curEntityId], models.UserReactionsWithinEntity{UserId: userId, Reactions: reactionIds})
	return sb
}

func (sb *stateBuilder) Build() state {
	var urs []models.UserReaction
	for entityId, usersReactions := range sb.entityId2usersReactions {
		e := sb.entityId2entity[entityId]
		for _, uRs := range usersReactions {
			for _, reactionId := range uRs.Reactions {
				ur := models.UserReaction{
					NamespaceId: e.GetNamespace().GetId(),
					EntityId:    e.GetId(),
					ReactionId:  reactionId,
					UserId:      uRs.UserId}
				urs = append(urs, ur)
			}
		}
	}
	return state{userReactions: urs}
}

type state struct {
	userReactions []models.UserReaction
}

func (s *state) GetReactionsCount() []ReactionCount {
	rcMap := make(map[string]int)
	for _, ur := range s.userReactions {
		rcMap[ur.ReactionId]++
	}
	var rc []ReactionCount
	for id, count := range rcMap {
		rc = append(rc, ReactionCount{ReactionId: id, Count: count})
	}
	return rc
}

func (s *state) GetMyReactionIds(userId string) []string {
	var res []string
	for _, ur := range s.userReactions {
		if ur.UserId == userId && !utils.Contains(res, ur.ReactionId) {
			res = append(res, ur.ReactionId)
		}
	}
	return res
}

func (s *state) Apply(e Entity, userId string) {
	s.ApplyToEntity(e, userId)
	s.ApplyToStorage()
}

func (s *state) ApplyToEntity(e Entity, userId string) {
	utils.AssertPtr(e)
	e.Update(s.GetReactionsCount(), s.GetMyReactionIds(userId))
}

func (s *state) ApplyToStorage() {
	setup.ClearUserReactions()
	for _, ur := range s.userReactions {
		var req models.ReactionsPOSTRequest
		force := false
		req.Query.Force = &force
		req.Body.NamespaceId = ur.NamespaceId
		req.Body.EntityId = ur.EntityId
		req.Body.ReactionId = ur.ReactionId
		req.Body.UserId = ur.UserId
		resp, err := setup.RSClient.AddReaction(&req)
		if err != nil {
			panic(fmt.Errorf("failed to add reaction: %w", err))
		}
		if resp.Code() != 200 {
			panic(fmt.Errorf("failed to add reaction: code != 200"))
		}
		setup.RS.RefreshEntityReactions(context.Background())
	}
}

func (s *state) RequireIsActual(t *testing.T, e Entity, userId string) {
	s.RequireIsActualForEntity(t, e, userId)
	s.RequireIsActualForStorage(t, e, userId)
}

func (s *state) RequireIsActualForEntity(t *testing.T, e Entity, userId string) {
	require.ElementsMatch(t, s.GetReactionsCount(), e.GetReactionsCount())
	require.ElementsMatch(t, s.GetMyReactionIds(userId), e.GetMyReactionIds())
}

func (s *state) RequireIsActualForStorage(t *testing.T, e Entity, userId string) {
	var req models.ReactionsGETRequest
	req.Query.NamespaceId = e.GetNamespace().GetId()
	req.Query.EntityId = e.GetId()
	req.Query.UserId = userId
	resp, err := setup.RSClient.GetReactions(&req)
	if err != nil {
		panic(fmt.Errorf("failed to get reactions: %w", err))
	}
	if resp.Code() != 200 {
		panic(fmt.Errorf("failed to get reactions: code != 200"))
	}
	resp200, ok := resp.(*models.ReactionsGETResponse200)
	if !ok {
		panic(fmt.Errorf("failed to convert to resp200"))
	}
	require.ElementsMatch(t, s.GetReactionsCount(), resp200.ReactionsCount)
	require.ElementsMatch(t, s.GetMyReactionIds(userId), resp200.UserReactions.Reactions)
}
