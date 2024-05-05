package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

var DebugMetrics = []*ginprometheus.Metric{
	GETReactions_AcquireWrap.Metric,
	GETReactions_GetEntityReactionsCountWrap.Metric,
	GETReactions_GetUniqEntityUserReactionsWrap.Metric,
	POSTReactions_AcquireWrap.Metric,
	POSTReactions_GetNamespaceWrap.Metric,
	POSTReactions_AddUserReactionWrap.Metric,
	DELETEReactions_AcquireWrap.Metric,
	DELETEReactions_RemoveUserReactionWrap.Metric,
}

func UnwrapDebugMetrics() {
	GETReactions_Acquire = GETReactions_AcquireWrap.Unwrap()
	GETReactions_GetEntityReactionsCount = GETReactions_GetEntityReactionsCountWrap.Unwrap()
	GETReactions_GetUniqEntityUserReactions = GETReactions_GetUniqEntityUserReactionsWrap.Unwrap()
	POSTReactions_Acquire = POSTReactions_AcquireWrap.Unwrap()
	POSTReactions_GetNamespace = POSTReactions_GetNamespaceWrap.Unwrap()
	POSTReactions_AddUserReaction = POSTReactions_AddUserReactionWrap.Unwrap()
	DELETEReactions_Acquire = DELETEReactions_AcquireWrap.Unwrap()
	DELETEReactions_RemoveUserReaction = DELETEReactions_RemoveUserReactionWrap.Unwrap()
}

var GETReactions_AcquireWrap = &MetricWrap[prometheus.Histogram]{
	&ginprometheus.Metric{
		ID:          "GET_reactions__acquire",
		Name:        "GET_reactions__acquire",
		Description: "GET_/reactions__Acquire duration in seconds.",
		Type:        "histogram",
	},
}

var GETReactions_Acquire prometheus.Histogram

var GETReactions_GetEntityReactionsCountWrap = &MetricWrap[prometheus.Histogram]{
	&ginprometheus.Metric{
		ID:          "GET_reactions__get_entity_reactions_count",
		Name:        "GET_reactions__get_entity_reactions_count",
		Description: "GET_/reactions__GetEntityReactionsCount duration in seconds.",
		Type:        "histogram",
	},
}

var GETReactions_GetEntityReactionsCount prometheus.Histogram

var GETReactions_GetUniqEntityUserReactionsWrap = &MetricWrap[prometheus.Histogram]{
	&ginprometheus.Metric{
		ID:          "GET_reactions__get_uniq_entity_user_reactions",
		Name:        "GET_reactions__get_uniq_entity_user_reactions",
		Description: "GET_/reactions__GetUniqEntityUserReactions duration in seconds.",
		Type:        "histogram",
	},
}

var GETReactions_GetUniqEntityUserReactions prometheus.Histogram

var POSTReactions_AcquireWrap = &MetricWrap[prometheus.Histogram]{
	&ginprometheus.Metric{
		ID:          "POST_reactions__acquire",
		Name:        "POST_reactions__acquire",
		Description: "POST_/reactions__Acquire duration in seconds.",
		Type:        "histogram",
	},
}

var POSTReactions_Acquire prometheus.Histogram

var POSTReactions_GetNamespaceWrap = &MetricWrap[prometheus.Histogram]{
	&ginprometheus.Metric{
		ID:          "POST_reactions__get_namespace",
		Name:        "POST_reactions__get_namespace",
		Description: "POST_/reactions__GetNamespace duration in seconds.",
		Type:        "histogram",
	},
}

var POSTReactions_GetNamespace prometheus.Histogram

var POSTReactions_AddUserReactionWrap = &MetricWrap[prometheus.Histogram]{
	&ginprometheus.Metric{
		ID:          "POST_reactions__add_user_reaction",
		Name:        "POST_reactions__add_user_reaction",
		Description: "POST_/reactions__AddUserReaction duration in seconds.",
		Type:        "histogram",
	},
}

var POSTReactions_AddUserReaction prometheus.Histogram

var DELETEReactions_AcquireWrap = &MetricWrap[prometheus.Histogram]{
	&ginprometheus.Metric{
		ID:          "DELETE_reactions__acquire",
		Name:        "DELETE_reactions__acquire",
		Description: "DELETE_/reactions__Acquire duration in seconds.",
		Type:        "histogram",
	},
}

var DELETEReactions_Acquire prometheus.Histogram

var DELETEReactions_RemoveUserReactionWrap = &MetricWrap[prometheus.Histogram]{
	&ginprometheus.Metric{
		ID:          "DELETE_reactions__remove_user_reaction",
		Name:        "DELETE_reactions__remove_user_reaction",
		Description: "DELETE_/reactions__RemoveUserReaction duration in seconds.",
		Type:        "histogram",
	},
}

var DELETEReactions_RemoveUserReaction prometheus.Histogram

func Record(fn func(), m prometheus.Histogram) {
	start := time.Now()
	fn()
	elapsed := float64(time.Since(start)) / float64(time.Second)
	if m != nil {
		m.Observe(elapsed)
	}
}
