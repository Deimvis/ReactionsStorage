package http_handlers_test

import (
	"testing"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/tests/fake"
)

func TestAvailableReactionsGETSimple(t *testing.T) {
	var req models.AvailableReactionsGETRequest
	req.Query.NamespaceId = fake.NamespaceId

	var resp models.AvailableReactionsGETResponse200
	resp.Reactions = []models.Reaction{fake.Reaction, fake.Reaction2, fake.Reaction3}

	test(t, &req, &resp)
}
