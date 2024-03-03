package http_handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/tests/fake"
)

func setFakeConfiguration() {
	ctx := context.Background()
	cs.ClearStrict(ctx)
	cs.AddReactionStrict(ctx, &fake.Reaction)
	cs.AddReactionStrict(ctx, &fake.Reaction2)
	cs.AddReactionStrict(ctx, &fake.Reaction3)
	cs.AddReactionSetStrict(ctx, &fake.ReactionSet)
	cs.AddNamespaceStrict(ctx, &fake.Namespace)
}

func clearUserReactions() {
	rs.ClearStrict(context.Background())
}

func test(t *testing.T, req models.Request, resp models.Response) {
	w := request(t, req)
	requireResponse(t, resp, w)
}

func request(t *testing.T, r models.Request) *httptest.ResponseRecorder {
	return requestRaw(t, r.Method(), fmt.Sprintf("%s?%s", r.Path(), r.QueryString()), bytes.NewReader(r.BodyRaw()))
}

func requestRaw(t *testing.T, method string, url string, body io.Reader) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, body)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	srv.Handler.ServeHTTP(w, req)
	return w
}

func requireResponse(t *testing.T, exp models.Response, act *httptest.ResponseRecorder) {
	expBody, err := json.Marshal(exp)
	require.NoError(t, err)
	requireResponseRaw(t, exp.Code(), string(expBody), act)
}

func requireResponseRaw(t *testing.T, expCode int, expBody string, act *httptest.ResponseRecorder) {
	require.Equal(t, expCode, act.Code)
	require.Equal(t, expBody, act.Body.String())
}
