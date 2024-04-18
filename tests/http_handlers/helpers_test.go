package http_handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Deimvis/reactionsstorage/src/models"
	setup "github.com/Deimvis/reactionsstorage/tests/setup"
)

func test(t *testing.T, req models.Request, resp models.Response) {
	w := request(t, req)
	requireResponse(t, resp, w)
}

func request(t *testing.T, r models.Request) *httptest.ResponseRecorder {
	return requestRaw(t, r.Method(), fmt.Sprintf("%s?%s", r.Path(), r.QueryString()), bytes.NewReader(r.BodyRaw()), r.Header())
}

func requestRaw(t *testing.T, method string, url string, body io.Reader, headers ...http.Header) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, body)
	require.NoError(t, err)
	req.Header = mergeHeaders(headers...)
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
	require.JSONEq(t, expBody, act.Body.String())
}

func mergeHeaders(headers ...http.Header) http.Header {
	result := make(http.Header)
	for _, headersGroup := range headers {
		for k, vs := range headersGroup {
			for _, v := range vs {
				result.Add(k, v)
			}
		}
	}
	return result
}
