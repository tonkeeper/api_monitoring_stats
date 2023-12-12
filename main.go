package main

import (
	"api_monitoring_stats/services/tonapi"
	"api_monitoring_stats/services/toncenter"
	"fmt"
	"log"
	"net/http"
	"time"

	"api_monitoring_stats/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config.LoadConfig()

	metricServer := http.Server{
		Addr:    fmt.Sprintf(":%v", config.Config.MetricsPort),
		Handler: promhttp.Handler(),
	}
	go func() {
		if err := metricServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start metrix server: %v", err)
		}
	}()
	sources := []metrics{
		tonapi.NewMonitoring(),
		toncenter.NewV2Monitoring("toncenter.com", "https://toncenter.com/api/v2/"),
	}
	go workerMetrics(sources)

	for {
		time.Sleep(time.Hour)
	}
}
