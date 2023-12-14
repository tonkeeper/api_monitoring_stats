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
	}
	go workerMetrics(sources)

	for {
		time.Sleep(time.Hour)
	}
}

//todo:
//curl 'https://api.dedust.io/v2/routing/plan'-X POST -H 'content-type: application/json'   --data-raw '{"from":"native","to":"jetton:0:65aac9b5e380eae928db3c8e238d9bc0d61a9320fdc2bc7a2f6c87d6fedf9208","amount":"1000000000"}'
//curl 'https://rpc.ston.fi/' --compressed -X POST -H 'Content-Type: application/json'  --data-raw '{"jsonrpc":"2.0","id":7,"method":"dex.simulate_swap","params":{"offer_address":"EQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAM9c","offer_units":"1000000000","ask_address":"EQA2kCVNwVsil2EM2mB0SkXytxCqQjS4mttjDpnXmwG9T6bO","slippage_tolerance":"0.001"}}'
