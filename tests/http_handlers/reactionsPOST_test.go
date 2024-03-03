package http_handlers_test

import (
	"testing"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/tests/fake"
)

func TestReactionsPOSTSimple(t *testing.T) {
	defer clearUserReactions()

	var req models.ReactionsPOSTRequest
	req.Body.NamespaceId = fake.NamespaceId
	req.Body.EntityId = fake.EntityId
	req.Body.ReactionId = fake.ReactionId
	req.Body.UserId = fake.UserId

	var resp models.ReactionsPOSTResponse200
	resp.Status = "ok"

	test(t, &req, &resp)
}
