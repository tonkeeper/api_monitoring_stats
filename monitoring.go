package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	MetricServiceTimeHistogramVec = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tonstatus_functions_time",
			Help:    "Service functions execution duration distribution in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.2, 0.3, 0.4, 0.5, 0.7, 1, 2.5, 5, 10},
		},
		[]string{"service"},
	)

	MetricServiceIndexingLatencyHistogramVec = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tonstatus_indexing_latency",
			Help:    "difference between current time and last transaction on electror",
			Buckets: []float64{0.5, 1, 2, 3, 4, 5, 7.5, 10, 15, 20, 30, 60, 120, 300, 1200, 3600},
		},
		[]string{"service"},
	)

	MetricServiceRequest = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tonstatus_http_api",
		Help: "availability of http api",
	}, []string{"service"})

	MetricTxPropagationLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tonstatus_tx_propagation_latency_seconds",
			Help:    "Time from tx sent to network until it appeared in the API",
			Buckets: []float64{0.5, 1, 2, 3, 5, 7, 10, 15, 20, 30, 45, 60, 90, 120},
		},
		[]string{"service"},
	)
)

var (
	MetricDAppTimeHistogramVec = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tonstatus_dapp_functions_time",
			Help:    "DApp functions execution duration distribution in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.2, 0.3, 0.4, 0.5, 0.7, 1, 2.5, 5, 10},
		},
		[]string{"dapp"})

	MetricDAppMainPageLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tonstatus_dapp_main_page_time",
			Buckets: []float64{0.01, 0.05, 0.1, 0.2, 0.3, 0.4, 0.5, 0.7, 1, 2.5, 5, 10},
		},
		[]string{"dapp"})

	MetricDAppIndexingLatencyHistogramVec = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tonstatus_dapp_indexing_latency",
			Help:    "difference between current time and last transaction on electror",
			Buckets: []float64{0.5, 1, 5, 7.5, 10, 12.5, 15, 17.5, 20, 25, 30, 60, 120, 300, 600, 1200, 3600},
		},
		[]string{"dapp"},
	)

	MetricDAppAvailability = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tonstatus_dapp_availability",
		Help: "availability of http api",
	}, []string{"dapp"})
)

var (
	metricBridgeAvailability = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tonconnect_bridge_availability",
	}, []string{"bridge"})
	metricBridgeLatencyHistogramVec = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "tonconnect_bridge_latency",
		Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2, 5, 10},
	}, []string{"bridge"})
	metricBridgeReconnects = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tonconnect_bridge_reconnects",
	}, []string{"bridge"})
)
