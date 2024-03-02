package http_handlers_test

import (
	"testing"

	"github.com/Deimvis/reactionsstorage/src/models"
)

func TestReactionsDELETESimple(t *testing.T) {
	defer clearUserReactions()

	var req models.ReactionsDELETERequest
	req.Body.NamespaceId = fakeNamespaceId
	req.Body.EntityId = fakeEntityId
	req.Body.ReactionId = fakeReactionId
	req.Body.UserId = fakeUserId

	var resp models.ReactionsDELETEResponse200
	resp.Status = "ok"

	test(t, &req, &resp)
}
