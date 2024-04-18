package metrics

import ginprometheus "github.com/zsais/go-gin-prometheus"

type MetricWrap[T any] struct {
	*ginprometheus.Metric
}

func (mw *MetricWrap[T]) Unwrap() T {
	return mw.MetricCollector.(T)
}
