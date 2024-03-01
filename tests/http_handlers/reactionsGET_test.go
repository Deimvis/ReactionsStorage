package http_handlers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReactionsGETSimple(t *testing.T) {
	w := requestRaw(t, "GET", "/reactions?namespace_id=namespace&entity_id=entity&user_id=user", nil)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"entity_id":"entity","reactions_count":null,"user_reactions":{"user_id":"user","reactions":[]}}`, w.Body.String())
}
