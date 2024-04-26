package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

var GETReactionsAcquireWrap = &MetricWrap[prometheus.Histogram]{
	&ginprometheus.Metric{
		ID:          "GET_reactions__acquire",
		Name:        "GET_reactions__acquire",
		Description: "GET_/reactions_Acquire.",
		Type:        "histogram",
	},
}

var GETReactionsAcquire prometheus.Histogram

var GetEntityReactionsCountWrap = &MetricWrap[prometheus.Histogram]{
	&ginprometheus.Metric{
		ID:          "get_entity_reactions_count",
		Name:        "get_entity_reactions_count",
		Description: "GetEntityReactionsCount.",
		Type:        "histogram",
	},
}

var GetEntityReactionsCount prometheus.Histogram

var GetUniqEntityUserReactionsWrap = &MetricWrap[prometheus.Histogram]{
	&ginprometheus.Metric{
		ID:          "get_uniq_entity_user_reactions",
		Name:        "get_uniq_entity_user_reactions",
		Description: "GetUniqEntityUserReactions.",
		Type:        "histogram",
	},
}

var GetUniqEntityUserReactions prometheus.Histogram

var GetEntityReactionsCountQueryWrap = &MetricWrap[prometheus.Histogram]{
	&ginprometheus.Metric{
		ID:          "get_uniq_entity_user_reactions__query",
		Name:        "get_uniq_entity_user_reactions__query",
		Description: "GetUniqEntityUserReactions_Query.",
		Type:        "histogram",
	},
}

var GetEntityReactionsCountQuery prometheus.Histogram

var GetEntityReactionsCountCollectRowsWrap = &MetricWrap[prometheus.Histogram]{
	&ginprometheus.Metric{
		ID:          "get_uniq_entity_user_reactions__collect_rows",
		Name:        "get_uniq_entity_user_reactions__collect_rows",
		Description: "GetUniqEntityUserReactions_CollectRows.",
		Type:        "histogram",
	},
}

var GetEntityReactionsCountCollectRows prometheus.Histogram

func Record(fn func(), m prometheus.Histogram) {
	start := time.Now()
	fn()
	elapsed := float64(time.Since(start)) / float64(time.Second)
	if m != nil {
		m.Observe(elapsed)
	}
}
