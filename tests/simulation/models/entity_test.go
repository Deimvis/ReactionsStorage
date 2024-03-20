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
		initialState  state
		addReactionId string
		expectOk      bool
		expectState   state
	}{
		{
			state{nil, nil},
			fake.ReactionId1,
			true,
			state{[]ReactionCount{{fake.ReactionId1, 1}}, []string{fake.ReactionId1}},
		},
		{
			state{[]ReactionCount{{fake.ReactionId1, 1}}, []string{fake.ReactionId1}},
			fake.ReactionId2,
			true,
			state{[]ReactionCount{{fake.ReactionId1, 1}, {fake.ReactionId2, 1}}, []string{fake.ReactionId1, fake.ReactionId2}},
		},
		{
			state{[]ReactionCount{{fake.ReactionId1, 1}}, []string{fake.ReactionId1}},
			fake.ReactionId1,
			false,
			state{[]ReactionCount{{fake.ReactionId1, 1}}, []string{fake.ReactionId1}},
		},
		{
			state{[]ReactionCount{{fake.ReactionId1, 1}, {fake.ReactionId2, 1}}, nil},
			fake.ReactionId1,
			true,
			state{[]ReactionCount{{fake.ReactionId1, 2}, {fake.ReactionId2, 1}}, []string{fake.ReactionId1}},
		},
		{
			state{[]ReactionCount{{fake.ReactionId1, 1}, {fake.ReactionId2, 1}}, []string{fake.ReactionId1}},
			fake.ReactionId2,
			true,
			state{[]ReactionCount{{fake.ReactionId1, 1}, {fake.ReactionId2, 2}}, []string{fake.ReactionId1, fake.ReactionId2}},
		},
	}

	n := NewNamespace(fake.NamespaceId, setup.RSClient)
	e := NewEntity("", n)
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			e.Update(tc.initialState.ReactionCount, tc.initialState.MyReactionIds)
			ok := e.AddMyReaction(tc.addReactionId)
			require.Equal(t, tc.expectOk, ok)
			require.ElementsMatch(t, tc.expectState.ReactionCount, e.GetReactionsCount())
			require.ElementsMatch(t, tc.expectState.MyReactionIds, e.GetMyReactionIds())
		})
	}
}

func TestEntity_RemoveMyReaction(t *testing.T) {
	testCases := []struct {
		initialState     state
		removeReactionId string
		expectOk         bool
		expectState      state
	}{
		{
			state{nil, nil},
			fake.ReactionId1,
			false,
			state{nil, nil},
		},
		{
			state{[]ReactionCount{{fake.ReactionId1, 100}, {fake.ReactionId2, 1}, {fake.ReactionId3, 1}}, []string{fake.ReactionId1, fake.ReactionId2}},
			fake.ReactionId1,
			true,
			state{[]ReactionCount{{fake.ReactionId1, 99}, {fake.ReactionId2, 1}, {fake.ReactionId3, 1}}, []string{fake.ReactionId2}},
		},
		{
			state{[]ReactionCount{{fake.ReactionId1, 100}, {fake.ReactionId2, 1}, {fake.ReactionId3, 1}}, []string{fake.ReactionId1, fake.ReactionId2}},
			fake.ReactionId2,
			true,
			state{[]ReactionCount{{fake.ReactionId1, 100}, {fake.ReactionId2, 0}, {fake.ReactionId3, 1}}, []string{fake.ReactionId1}},
		},
		{
			state{[]ReactionCount{{fake.ReactionId1, 100}, {fake.ReactionId2, 0}, {fake.ReactionId3, 1}}, []string{fake.ReactionId1}},
			fake.ReactionId2,
			false,
			state{[]ReactionCount{{fake.ReactionId1, 100}, {fake.ReactionId2, 0}, {fake.ReactionId3, 1}}, []string{fake.ReactionId1}},
		},
	}

	n := NewNamespace(fake.NamespaceId, setup.RSClient)
	e := NewEntity("", n)
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			e.Update(tc.initialState.ReactionCount, tc.initialState.MyReactionIds)
			ok := e.RemoveMyReaction(tc.removeReactionId)
			require.Equal(t, tc.expectOk, ok)
			require.ElementsMatch(t, tc.expectState.ReactionCount, e.GetReactionsCount())
			require.ElementsMatch(t, tc.expectState.MyReactionIds, e.GetMyReactionIds())
		})
	}
}

type state struct {
	ReactionCount []ReactionCount
	MyReactionIds []string
}
