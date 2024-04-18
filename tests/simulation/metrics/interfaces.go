package metrics

import "net/http"

type HTTPHandler = func(r *http.Request) (*http.Response, error)

type HTTPRecorder interface {
	Record(handler HTTPHandler, req *http.Request) (*http.Response, error)
	Sync()
}
