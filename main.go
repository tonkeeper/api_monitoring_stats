package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"api_monitoring_stats/config"
	"api_monitoring_stats/services/dton"
	public_config "api_monitoring_stats/services/public-config"
	"api_monitoring_stats/services/tonapi"
	"api_monitoring_stats/services/toncenter"
	"api_monitoring_stats/services/tonhub"
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
		dton.NewMonitoring("dton", "https://dton.io/graphql"),
		toncenter.NewV2Monitoring("toncenter.com v2", "https://toncenter.com/api/v2"),
		toncenter.NewV3Monitoring("toncenter.com v3", "https://toncenter.com/api/v3"),
		toncenter.NewV2Monitoring("orbs http-api", "https://ton.access.orbs.network/route/1/mainnet/toncenter-api-v2"),
		tonhub.NewV4Monitoring("tonhub", "https://mainnet-v4.tonhubapi.com"),
		public_config.NewLiteServersMetrics(),
		dton.NewMonitoring("dton", "https://dton.io/graphql"),
	}
	go workerMetrics(sources)

	for {
		time.Sleep(time.Hour)
	}
}
