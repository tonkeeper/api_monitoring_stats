package main

import (
	"time"

	"api_monitoring_stats/services"
	"api_monitoring_stats/services/tonapi"
)

func workerMetrics() {
	for {
		tonApiMetrics, _ := tonapi.NewMonitoring().GetMetrics()

		metrics := []services.ApiMetrics{tonApiMetrics}
		for _, metric := range metrics {
			serviceName := string(metric.ServiceName)
			services.MetricServiceTimeHistogramVec.WithLabelValues(serviceName).Observe(metric.HttpsLatency)
			services.MetricServiceIndexingLatencyHistogramVec.WithLabelValues(serviceName).Observe(metric.IndexingLatency)
			if metric.Alive {
				services.MetricServiceRequestSuccess.WithLabelValues(serviceName).Inc()
			} else {
				services.MetricServiceRequestFails.WithLabelValues(serviceName).Inc()
			}
		}

		time.Sleep(1 * time.Minute)
	}
}
