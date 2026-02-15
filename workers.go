package main

import (
	"context"
	"fmt"
	"time"

	"api_monitoring_stats/services"
)

type metrics[T services.DAppMetrics | services.ApiMetrics | services.BridgeMetrics | services.TxPropagationMetrics] interface {
	GetMetrics(ctx context.Context) T
}

func workerMetrics[T services.ApiMetrics | services.DAppMetrics | services.BridgeMetrics | services.TxPropagationMetrics](sources []metrics[T], f func(m T)) {
	time.Sleep(time.Second)
	for _, s := range sources {
		go func(m metrics[T]) {
			for {
				collect(m, f)
				sleep := time.Second * 30
				if i, ok := m.(interface{ CheckInterval() time.Duration }); ok {
					sleep = i.CheckInterval()
				}
				time.Sleep(sleep)
			}
		}(s)
		time.Sleep(time.Second * 2)
	}
}

func collect[T services.ApiMetrics | services.DAppMetrics | services.BridgeMetrics | services.TxPropagationMetrics](s metrics[T], f func(m T)) {
	timeout := 10 * time.Second
	if t, ok := s.(interface{ Timeout() time.Duration }); ok {
		timeout = t.Timeout()
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	m := s.GetMetrics(ctx)
	f(m)
}

func apiMetricsCollect(m services.ApiMetrics) {
	MetricServiceTimeHistogramVec.WithLabelValues(m.ServiceName).Observe(m.HttpsLatency)
	MetricServiceIndexingLatencyHistogramVec.WithLabelValues(m.ServiceName).Observe(m.IndexingLatency)
	MetricServiceRequest.WithLabelValues(m.ServiceName).Set(float64(m.SuccessChecks) / float64(m.TotalChecks))
	for _, err := range m.Errors {
		fmt.Println("Service", m.ServiceName, err.Error())
	}
}

func dappsMetricsCollect(m services.DAppMetrics) {
	MetricDAppAvailability.WithLabelValues(m.ServiceName).Set(float64(m.SuccessChecks) / float64(m.TotalChecks))
	MetricDAppMainPageLatency.WithLabelValues(m.ServiceName).Observe(m.MainPageLoadLatency)
	MetricDAppTimeHistogramVec.WithLabelValues(m.ServiceName).Observe(m.ApiLatency)
	if m.IndexationLatency != nil {
		MetricDAppIndexingLatencyHistogramVec.WithLabelValues(m.ServiceName).Observe(*m.IndexationLatency)
	}
	for _, err := range m.Errors {
		fmt.Println("Service", m.ServiceName, err.Error())
	}
}

func bridgeMetricsCollect(m services.BridgeMetrics) {
	metricBridgeAvailability.WithLabelValues(m.ServiceName).Set(float64(m.SuccessChecks) / float64(m.TotalChecks))
	metricBridgeReconnects.WithLabelValues(m.ServiceName).Set(float64(m.Reconnects))
	metricBridgeLatencyHistogramVec.WithLabelValues(m.ServiceName).Observe(m.TransferLatency)
	for _, err := range m.Errors {
		fmt.Println("Service", m.ServiceName, err.Error())
	}
}

func txPropagationCollect(m services.TxPropagationMetrics) {
	MetricTxPropagationLatency.WithLabelValues(m.ServiceName).Observe(m.Latency)
	for _, err := range m.Errors {
		fmt.Println("Service", m.ServiceName, err.Error())
	}
}
