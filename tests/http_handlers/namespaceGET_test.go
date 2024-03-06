package http_handlers_test

import (
	"testing"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/tests/fake"
)

func TestNamespaceGETSimple(t *testing.T) {
	var req models.NamespaceGETRequest
	req.Query.NamespaceId = fake.NamespaceId

	var resp models.NamespaceGETResponse200
	resp.Namespace = fake.Namespace

	test(t, &req, &resp)
}
