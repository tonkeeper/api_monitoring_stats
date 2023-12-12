package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const Dead = "dead"
const Alive = "alive"
const Undead = "degraded"

var (
	MetricServiceTimeHistogramVec = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "service_functions_time",
			Help:    "Service functions execution duration distribution in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 1, 5, 10},
		},
		[]string{"service"},
	)

	MetricServiceIndexingLatencyHistogramVec = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "service_indexing_latency",
			Help:    "",
			Buckets: []float64{1, 5, 10, 30, 60, 120, 300, 600, 1200, 1800, 3600},
		},
		[]string{"service"},
	)

	MetricServiceRequest = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "service_requests",
		Help: "The total number of success request",
	}, []string{"service", "status"})
)
