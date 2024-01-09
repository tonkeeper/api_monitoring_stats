package dapps

import (
	"context"
	"strings"

	"api_monitoring_stats/services"
)

type DeDust struct {
	dAppUrl    string
	calcApiUrl string
}

func NewDeDust() *DeDust {
	return &DeDust{
		dAppUrl:    "https://dedust.io/swap",
		calcApiUrl: "https://api.dedust.io/v2/routing/plan",
	}
}

func (d *DeDust) GetMetrics(ctx context.Context) services.DAppMetrics {
	m := services.DAppMetrics{
		ServiceName: "DeDust",
	}
	m.MainPageLoadLatency = services.HttpGet(&m.TotalChecks, &m.SuccessChecks, &m.Errors, d.dAppUrl, nil)
	m.ApiLatency = services.HttpPost(&m.TotalChecks, &m.SuccessChecks, &m.Errors, d.calcApiUrl, strings.NewReader(`{"from":"native","to":"jetton:0:65aac9b5e380eae928db3c8e238d9bc0d61a9320fdc2bc7a2f6c87d6fedf9208","amount":"1000000000"}`), nil)
	return m
}
