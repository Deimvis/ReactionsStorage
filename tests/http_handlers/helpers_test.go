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
	setup "github.com/Deimvis/reactionsstorage/tests/setup"
)

func setFakeConfiguration() {
	ctx := context.Background()
	setup.CS.ClearStrict(ctx)
	setup.CS.AddReactionStrict(ctx, &fake.Reaction)
	setup.CS.AddReactionStrict(ctx, &fake.Reaction2)
	setup.CS.AddReactionStrict(ctx, &fake.Reaction3)
	setup.CS.AddReactionSetStrict(ctx, &fake.ReactionSet)
	setup.CS.AddNamespaceStrict(ctx, &fake.Namespace)
}

func clearUserReactions() {
	setup.RS.ClearStrict(context.Background())
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
	setup.SRV.Handler.ServeHTTP(w, req)
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
