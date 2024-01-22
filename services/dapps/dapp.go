package dapps

import (
	"context"
	"strings"

	"api_monitoring_stats/services"
)

type Dapp struct {
	name        string
	dAppUrl     string
	calcApiUrl  string
	calcPayload *string
}

func (d *Dapp) GetMetrics(ctx context.Context) services.DAppMetrics {
	m := services.DAppMetrics{
		ServiceName: d.name,
	}
	m.MainPageLoadLatency = services.HttpGet(ctx, &m.TotalChecks, &m.SuccessChecks, &m.Errors, d.dAppUrl, nil)
	if d.calcApiUrl != "" {
		if d.calcPayload != nil {
			m.ApiLatency = services.HttpPost(ctx, &m.TotalChecks, &m.SuccessChecks, &m.Errors, d.calcApiUrl, strings.NewReader(*d.calcPayload), nil)
		} else {
			m.ApiLatency = services.HttpGet(ctx, &m.TotalChecks, &m.SuccessChecks, &m.Errors, d.calcApiUrl, nil)
		}
	}
	return m
}

func pointer[T any](s T) *T {
	return &s
}
