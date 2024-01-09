package dapps

import (
	"context"
	"strings"

	"api_monitoring_stats/services"
)

type StonFi struct {
	dAppUrl    string
	calcApiUrl string
}

func NewStonFi() *StonFi {
	return &StonFi{
		dAppUrl:    "https://app.ston.fi/swap?ft=jUSDT&tt=STON",
		calcApiUrl: "https://rpc.ston.fi/",
	}
}

func (d *StonFi) GetMetrics(ctx context.Context) services.DAppMetrics {
	m := services.DAppMetrics{
		ServiceName: "StonFi",
	}
	m.MainPageLoadLatency = services.HttpGet(&m.TotalChecks, &m.SuccessChecks, &m.Errors, d.dAppUrl, nil)
	m.ApiLatency = services.HttpPost(&m.TotalChecks, &m.SuccessChecks, &m.Errors, d.calcApiUrl, strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"dex.simulate_swap","params":{"offer_address":"EQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAM9c","offer_units":"1000000000","ask_address":"EQA2kCVNwVsil2EM2mB0SkXytxCqQjS4mttjDpnXmwG9T6bO","slippage_tolerance":"0.001"}}`), nil)
	return m
}
