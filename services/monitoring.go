package services

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type serviceName string

const (
	TonApi    serviceName = "TonAPI"
	DTon      serviceName = "dTON"
	TonCenter serviceName = "TonCenter"
)

var (
	MetricServiceTimeHistogramVec = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "service_functions_time",
			Help:    "Service functions execution duration distribution in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 1, 5, 10, 60},
		},
		[]string{"service"},
	)

	MetricServiceIndexingLatencyHistogramVec = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "service_indexing_latency",
			Help:    "", // TODO: write great help
			Buckets: []float64{0.01, 0.05, 0.1, 1, 5, 10, 60},
		},
		[]string{"service"},
	)

	MetricServiceRequestSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "service_request_success",
		Help: "The total number of success request",
	}, []string{"service"})

	MetricServiceRequestFails = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "service_request_fails",
		Help: "The total number of failed request",
	}, []string{"service"})
)

type ApiMetrics struct {
	ServiceName     serviceName
	HttpsLatency    float64
	Alive           bool
	IndexingLatency float64
}

type metrics interface {
	GetMetrics() (ApiMetrics, error)
}
