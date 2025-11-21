package main

import (
	"api_monitoring_stats/services"
	"api_monitoring_stats/services/connect"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var defaultBridges = []metrics[services.BridgeMetrics]{
	connect.NewBridge("tonapi", "https://bridge.tonapi.io/bridge"),
	connect.NewBridge("MTW", "https://tonconnectbridge.mytonwallet.org/bridge"),
	connect.NewBridge("tonhub", "https://connect.tonhubapi.com/tonconnect"),
	connect.NewBridge("TonSpace", "https://bridge.ton.space/bridge"),
	connect.NewBridge("DeWallet", "https://bridge.dewallet.pro/bridge"),
}

func GetBridgeConfig(configUrl string) ([]metrics[services.BridgeMetrics], error) {
	if configUrl == "" {
		return defaultBridges, nil
	}

	bridges := []metrics[services.BridgeMetrics]{}
	req, err := http.NewRequest(http.MethodGet, configUrl, nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return defaultBridges, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return defaultBridges, fmt.Errorf("failed to get wallets config: %v", err)
	}
	defer resp.Body.Close()

	type WalletConfig struct {
		AppName string `json:"app_name"`
		Bridges []struct {
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"bridge"`
	}
	var config []WalletConfig
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return defaultBridges, fmt.Errorf("failed to read wallets config: %v", err)
	}
	err = json.NewDecoder(bytes.NewReader(respBody)).Decode(&config)
	if err != nil {
		return defaultBridges, fmt.Errorf("failed to decode wallets config: %v", err)
	}

	for _, wallet := range config {
		for _, bridge := range wallet.Bridges {
			if bridge.Type == "sse" {
				bridges = append(bridges, connect.NewBridge(wallet.AppName, bridge.URL))
			}
		}
	}

	if len(bridges) == 0 {
		return defaultBridges, fmt.Errorf("no bridges found")
	}
	return bridges, nil
}
