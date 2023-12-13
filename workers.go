package main

import (
	"context"
	"fmt"
	"time"

	"api_monitoring_stats/services"
)

type metrics interface {
	GetMetrics(ctx context.Context) services.ApiMetrics
}

func workerMetrics(sources []metrics) {
	for _, s := range sources {
		go func(m metrics) {
			for {
				collect(m)
				time.Sleep(time.Second * 30)
			}
		}(s)
		time.Sleep(time.Second * 2)
	}
}

func collect(s metrics) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	m := s.GetMetrics(ctx)
	MetricServiceTimeHistogramVec.WithLabelValues(m.ServiceName).Observe(m.HttpsLatency)
	MetricServiceIndexingLatencyHistogramVec.WithLabelValues(m.ServiceName).Observe(m.IndexingLatency)

	MetricServiceRequest.WithLabelValues(m.ServiceName).Set(float64(m.SuccessChecks) / float64(m.TotalChecks))

	for _, err := range m.Errors {
		fmt.Println("Service", m.ServiceName, err.Error())
	}
}
