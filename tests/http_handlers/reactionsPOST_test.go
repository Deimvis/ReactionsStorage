package http_handlers_test

import (
	"testing"

	"github.com/Deimvis/reactionsstorage/src/models"
)

func TestReactionsPOSTSimple(t *testing.T) {
	defer clearUserReactions()

	var req models.ReactionsPOSTRequest
	req.Body.NamespaceId = fakeNamespaceId
	req.Body.EntityId = fakeEntityId
	req.Body.ReactionId = fakeReactionId
	req.Body.UserId = fakeUserId

	var resp models.ReactionsPOSTResponse200
	resp.Status = "ok"

	test(t, &req, &resp)
}
