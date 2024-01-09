package main

import (
	"time"

	"api_monitoring_stats/services"
)

type metricsWithPeriod[T services.ApiMetrics | services.DAppMetrics] struct {
	metrics[T]
	interval time.Duration
}

func Period[T services.ApiMetrics | services.DAppMetrics](m metrics[T], interval time.Duration) metrics[T] {
	return metricsWithPeriod[T]{
		metrics:  m,
		interval: interval,
	}
}

func (m metricsWithPeriod[T]) CheckInterval() time.Duration {
	return m.interval
}
