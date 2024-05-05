package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

var GinMetrics = []*ginprometheus.Metric{
	GINReqDurV2Wrap.Metric,
}

func UnwrapGinMetrics() {
	GINReqDurV2 = GINReqDurV2Wrap.Unwrap()
}

var GINReqDurV2Wrap = &MetricWrap[*prometheus.HistogramVec]{
	&ginprometheus.Metric{
		ID:          "GINReqDurV2",
		Name:        "request_duration_seconds_v2",
		Description: "The HTTP request latencies in seconds (v2).",
		Type:        "histogram_vec",
		Args:        []string{"code", "method", "url"},
	},
}

var GINReqDurV2 *prometheus.HistogramVec
