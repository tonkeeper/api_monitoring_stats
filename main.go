package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"api_monitoring_stats/config"
	"api_monitoring_stats/services"
	"api_monitoring_stats/services/connect"
	"api_monitoring_stats/services/dapps"
	"api_monitoring_stats/services/dton"
	public_config "api_monitoring_stats/services/public-config"
	"api_monitoring_stats/services/tonapi"
	"api_monitoring_stats/services/toncenter"
	"api_monitoring_stats/services/tonhub"
	"api_monitoring_stats/services/tonxapi"

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
	apis := []metrics[services.ApiMetrics]{
		tonapi.NewMonitoring(),
		dton.NewMonitoring("dton", "https://dton.co/"+config.Config.DtonToken+"/graphql_private"),
		toncenter.NewV2Monitoring("toncenter.com v2", "https://toncenter.com/api/v2", config.Config.TonCenterApiToken),
		toncenter.NewV3Monitoring("toncenter.com v3", "https://toncenter.com/api/v3", config.Config.TonCenterApiToken),
		toncenter.NewV2Monitoring("orbs http-api", "https://ton.access.orbs.network/route/1/mainnet/toncenter-api-v2", ""),
		tonhub.NewV4Monitoring("tonhub", "https://mainnet-v4.tonhubapi.com"),
		public_config.NewLiteServersMetrics("public liteservers", nil),
		tonxapi.NewTonXAPIMonitoring(
			"TonXAPI",
			"https://mainnet-rpc.tonxapi.com/v2/json-rpc",
		),
	}
	if config.Config.ChainstackToken != "" {
		apis = append(apis, Period[services.ApiMetrics](toncenter.NewV2Monitoring("chainstack.com v2", "https://ton-mainnet.core.chainstack.com/"+config.Config.ChainstackToken+"/api/v2", ""), time.Minute))
		apis = append(apis, Period[services.ApiMetrics](toncenter.NewV3Monitoring("chainstack.com v3", "https://ton-mainnet.core.chainstack.com/"+config.Config.ChainstackToken+"/api/v3", ""), time.Minute))
	}
	if config.Config.GetBlockKey != "" {
		apis = append(apis, Period[services.ApiMetrics](toncenter.NewV2Monitoring("getblock.io", "https://go.getblock.io/"+config.Config.GetBlockKey, ""), time.Minute))
	}
	if config.Config.DtonLiteServers != nil {
		apis = append(apis, public_config.NewLiteServersMetrics("liteservers_bot", config.Config.DtonLiteServers))
	}

	dappsMetrics := []metrics[services.DAppMetrics]{
		dapps.DeDust,
		dapps.StonFi,
		dapps.Getgems,
	}
	bridges := []metrics[services.BridgeMetrics]{
		connect.NewBridge("tonapi", "https://bridge.tonapi.io/bridge"),
		connect.NewBridge("MTW", "https://tonconnectbridge.mytonwallet.org/bridge"),
		connect.NewBridge("tonhub", "https://connect.tonhubapi.com/tonconnect"),
		connect.NewBridge("TonSpace", "https://bridge.ton.space/bridge"),
		connect.NewBridge("DeWallet", "https://bridge.dewallet.pro/bridge"),
	}
	go workerMetrics(apis, apiMetricsCollect)
	go workerMetrics(dappsMetrics, dappsMetricsCollect)
	go workerMetrics(bridges, bridgeMetricsCollect)
	for {
		time.Sleep(time.Hour)
	}
}
