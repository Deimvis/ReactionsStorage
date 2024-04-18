package metrics

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type PrometheusPushgatewayRecorder struct {
	pusher *push.Pusher

	reqCnt *prometheus.CounterVec
	reqDur *prometheus.HistogramVec
}

func NewPrometheusPushgatewayRecorder(host string, port int, ssl bool, subsystem string) *PrometheusPushgatewayRecorder {
	baseUrl := &url.URL{}
	if ssl {
		baseUrl.Scheme = "https"
	} else {
		baseUrl.Scheme = "http"
	}
	baseUrl.Host = fmt.Sprintf("%s:%d", host, port)

	reqCnt := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: subsystem,
			Name:      "requests_total",
			Help:      "Requests counter",
		},
		[]string{"code", "method", "host", "path", "url", "error"},
	)
	reqDur := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Subsystem: subsystem,
			Name:      "request_duration_seconds",
			Help:      "HTTP request latencies in seconds",
		},
		[]string{"code", "method", "host", "path", "url", "error"},
	)
	pusher := push.New(baseUrl.String(), subsystem).
		Collector(reqCnt).
		Collector(reqDur)

	return &PrometheusPushgatewayRecorder{
		pusher: pusher,
		reqCnt: reqCnt,
		reqDur: reqDur,
	}
}

func (r *PrometheusPushgatewayRecorder) Record(handler HTTPHandler, req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := handler(req)
	end := time.Now()

	code := ""
	if err == nil {
		code = strconv.Itoa(resp.StatusCode)
	}
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	elapsed := float64(end.Sub(start)) / float64(time.Second)
	r.reqDur.WithLabelValues(code, req.Method, req.Host, req.URL.Path, req.URL.String(), errMsg).Observe(elapsed)
	r.reqCnt.WithLabelValues(code, req.Method, req.Host, req.URL.Path, req.URL.String(), errMsg).Inc()

	return resp, err
}

func (r *PrometheusPushgatewayRecorder) Sync() {
	r.pusher.Push()
}
