package models

import (
	"fmt"
	"testing"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/tests/fake"
	"github.com/Deimvis/reactionsstorage/tests/setup"
)

func TestApp_Refresh(t *testing.T) {
	testCases := []struct {
		initialState state
		reactions    []models.UserReaction
		expectOk     bool
		expectState  state
	}{}

	n := NewNamespace(fake.NamespaceId, setup.RSClient)
	app := NewApp(setup.RSClient, []Topic{NewTopic("", 1, n)}, 1)
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			setup.ClearUserReactions()
			for _, r := range tc.reactions {
				var req models.ReactionsPOSTRequest
				req.Body.NamespaceId = r.NamespaceId
				req.Body.EntityId = r.EntityId
				req.Body.ReactionId = r.ReactionId
				req.Body.UserId = r.UserId
				setup.RSClient.AddReaction(&req)
			}
		})
	}
	setup.ClearUserReactions()
}

// func TestApp_AddReaction(t *testing.T) {
// 	testCases := []struct {
// 		initialState     state
// 		removeReactionId string
// 		expectOk         bool
// 		expectState      state
// 	}{

// 	}
// }
