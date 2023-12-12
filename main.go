package main

import (
	"fmt"
	"log"
	"net/http"

	"api_monitoring_stats/config"
	"api_monitoring_stats/controllers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tonkeeper/opentonapi/tonapi"
)

func main() {
	config.LoadConfig()

	tonApiClient, err := tonapi.New(tonapi.WithToken(config.Config.TonApiToken))
	if err != nil {
		log.Fatalf("[NewHandler] failed to init tonapi client: %v", err)
	}

	controllers.NewMetric(tonApiClient)

	handlers, err := controllers.NewHandler()
	if err != nil {
		log.Fatalf("failed to create api handler: %v", err)
	}

	metricServer := http.Server{
		Addr:    fmt.Sprintf(":%v", config.Config.MetricsPort),
		Handler: promhttp.Handler(),
	}
	go func() {
		if err = metricServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start metrix server: %v", err)
		}
	}()

	port := fmt.Sprintf(":%d", config.Config.Port)
	server, err := controllers.NewServer(handlers, port)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	log.Printf("Starting server on %v port", port)
	server.Run()
}
