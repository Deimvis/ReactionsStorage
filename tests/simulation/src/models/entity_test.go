package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Deimvis/reactionsstorage/tests/fake"
	"github.com/Deimvis/reactionsstorage/tests/setup"
)

func TestNewEntity(t *testing.T) {
	n := NewNamespace(fake.NamespaceId, setup.RSClient)
	e := NewEntity("", n)
	require.True(t, len(e.GetReactionsCount()) == 0)
	require.True(t, len(e.GetMyReactionIds()) == 0)
}

func TestEntity_AddMyReaction(t *testing.T) {
	testCases := []struct {
		initial       Entity
		addReactionId string
		expectOk      bool
		expect        Entity
	}{
		{
			makeEntity(nil, nil),
			fake.ReactionId1,
			true,
			makeEntity([]ReactionCount{{fake.ReactionId1, 1}}, []string{fake.ReactionId1}),
		},
		{
			makeEntity([]ReactionCount{{fake.ReactionId1, 1}}, []string{fake.ReactionId1}),
			fake.ReactionId2,
			true,
			makeEntity([]ReactionCount{{fake.ReactionId1, 1}, {fake.ReactionId2, 1}}, []string{fake.ReactionId1, fake.ReactionId2}),
		},
		{
			makeEntity([]ReactionCount{{fake.ReactionId1, 1}}, []string{fake.ReactionId1}),
			fake.ReactionId1,
			false,
			makeEntity([]ReactionCount{{fake.ReactionId1, 1}}, []string{fake.ReactionId1}),
		},
		{
			makeEntity([]ReactionCount{{fake.ReactionId1, 1}, {fake.ReactionId2, 1}}, nil),
			fake.ReactionId1,
			true,
			makeEntity([]ReactionCount{{fake.ReactionId1, 2}, {fake.ReactionId2, 1}}, []string{fake.ReactionId1}),
		},
		{
			makeEntity([]ReactionCount{{fake.ReactionId1, 1}, {fake.ReactionId2, 1}}, []string{fake.ReactionId1}),
			fake.ReactionId2,
			true,
			makeEntity([]ReactionCount{{fake.ReactionId1, 1}, {fake.ReactionId2, 2}}, []string{fake.ReactionId1, fake.ReactionId2}),
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			ok := tc.initial.AddMyReaction(tc.addReactionId)
			require.Equal(t, tc.expectOk, ok)
			require.ElementsMatch(t, tc.initial.GetReactionsCount(), tc.expect.GetReactionsCount())
			require.ElementsMatch(t, tc.initial.GetMyReactionIds(), tc.expect.GetMyReactionIds())
		})
	}
}

func TestEntity_RemoveMyReaction(t *testing.T) {
	makeEntity := func(rcs []ReactionCount, myReactionIds []string) Entity {
		n := NewNamespace(fake.NamespaceId, setup.RSClient)
		e := NewEntity("", n)
		e.Update(rcs, myReactionIds)
		return e
	}

	testCases := []struct {
		initial          Entity
		removeReactionId string
		expectOk         bool
		expect           Entity
	}{
		{
			makeEntity(nil, nil),
			fake.ReactionId1,
			false,
			makeEntity(nil, nil),
		},
		{
			makeEntity([]ReactionCount{{fake.ReactionId1, 100}, {fake.ReactionId2, 1}, {fake.ReactionId3, 1}}, []string{fake.ReactionId1, fake.ReactionId2}),
			fake.ReactionId1,
			true,
			makeEntity([]ReactionCount{{fake.ReactionId1, 99}, {fake.ReactionId2, 1}, {fake.ReactionId3, 1}}, []string{fake.ReactionId2}),
		},
		{
			makeEntity([]ReactionCount{{fake.ReactionId1, 100}, {fake.ReactionId2, 1}, {fake.ReactionId3, 1}}, []string{fake.ReactionId1, fake.ReactionId2}),
			fake.ReactionId2,
			true,
			makeEntity([]ReactionCount{{fake.ReactionId1, 100}, {fake.ReactionId2, 0}, {fake.ReactionId3, 1}}, []string{fake.ReactionId1}),
		},
		{
			makeEntity([]ReactionCount{{fake.ReactionId1, 100}, {fake.ReactionId2, 0}, {fake.ReactionId3, 1}}, []string{fake.ReactionId1}),
			fake.ReactionId2,
			false,
			makeEntity([]ReactionCount{{fake.ReactionId1, 100}, {fake.ReactionId2, 0}, {fake.ReactionId3, 1}}, []string{fake.ReactionId1}),
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			ok := tc.initial.RemoveMyReaction(tc.removeReactionId)
			require.Equal(t, tc.expectOk, ok)
			require.ElementsMatch(t, tc.initial.GetReactionsCount(), tc.expect.GetReactionsCount())
			require.ElementsMatch(t, tc.initial.GetMyReactionIds(), tc.expect.GetMyReactionIds())
		})
	}
}

func makeEntity(rcs []ReactionCount, myReactionIds []string) Entity {
	n := NewNamespace(fake.NamespaceId, setup.RSClient)
	e := NewEntity("", n)
	e.Update(rcs, myReactionIds)
	return e
}
