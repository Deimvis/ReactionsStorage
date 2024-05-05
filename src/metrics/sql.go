package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

var SQLMetrics = []*ginprometheus.Metric{
	SQLReqCntWrap.Metric,
	SQLReqDurWrap.Metric,
}

func UnwrapSQLMetrics() {
	SQLReqCnt = SQLReqCntWrap.Unwrap()
	SQLReqDur = SQLReqDurWrap.Unwrap()
}

var SQLReqCntWrap = &MetricWrap[prometheus.Counter]{
	&ginprometheus.Metric{
		ID:          "SQLReqCnt",
		Name:        "sql_query_count",
		Description: "The SQL query total count.",
		Type:        "counter",
	},
}

var SQLReqCnt prometheus.Counter

var SQLReqDurWrap = &MetricWrap[prometheus.Histogram]{
	&ginprometheus.Metric{
		ID:          "SQLReqDur",
		Name:        "sql_query_duration_seconds",
		Description: "The SQL query latencies in seconds..",
		Type:        "histogram",
	},
}

var SQLReqDur prometheus.Histogram
