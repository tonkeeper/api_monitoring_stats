package main

import (
	"api_monitoring_stats/services"
	"context"
	"fmt"
	"time"
)

type metrics interface {
	GetMetrics(ctx context.Context) services.ApiMetrics
}

func workerMetrics(sources []metrics) {
	for {
		for _, s := range sources {
			collect(s)
		}
		time.Sleep(time.Minute)
	}
}

func collect(s metrics) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	m := s.GetMetrics(ctx)
	MetricServiceTimeHistogramVec.WithLabelValues(m.ServiceName).Observe(m.HttpsLatency)
	MetricServiceIndexingLatencyHistogramVec.WithLabelValues(m.ServiceName).Observe(m.IndexingLatency)
	status := Alive
	if m.SuccessChecks == 0 {
		status = Dead
	} else if m.SuccessChecks < m.TotalChecks {
		status = Undead
	}
	MetricServiceRequest.WithLabelValues(m.ServiceName, status).Inc()

	for _, err := range m.Errors {
		fmt.Println("Service", m.ServiceName, err.Error())
	}
}
