package main

import (
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

	go workerMetrics()

	for {
		time.Sleep(time.Hour)
	}
}
