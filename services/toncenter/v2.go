package toncenter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"api_monitoring_stats/config"
	"api_monitoring_stats/services"
)

type V2Monitoring struct {
	name   string
	prefix string
}

func NewV2Monitoring(name, prefix string) *V2Monitoring {
	return &V2Monitoring{
		name:   name,
		prefix: prefix,
	}
}

func (vm *V2Monitoring) GetMetrics(ctx context.Context) services.ApiMetrics {
	m := services.ApiMetrics{
		ServiceName: vm.name,
	}
	m.TotalChecks++
	t := time.Now()
	r, err := http.Get(fmt.Sprintf("%v/getAddressInformation?address=%v", vm.prefix, config.ElectorAccountID.ToHuman(true, false)))
	if err != nil || r.StatusCode != http.StatusOK {
		m.Errors = append(m.Errors, fmt.Errorf("failed to get account state: %w", err))
	} else {
		m.SuccessChecks++
	}
	m.HttpsLatency = time.Since(t).Seconds()
	m.TotalChecks++

	time.Sleep(1 * time.Second)
	r, err = http.Get(fmt.Sprintf("%v/getTransactions?address=%v&limit=1&to_lt=0&archival=false", vm.prefix, config.ElectorAccountID.ToHuman(true, false)))
	if err != nil {
		m.Errors = append(m.Errors, fmt.Errorf("failed to get account transactions: %w", err))
		return m
	}
	defer r.Body.Close()
	var body struct {
		Result []struct {
			Utime int64 `json:"utime"`
		} `json:"result"`
	}
	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		m.Errors = append(m.Errors, fmt.Errorf("failed to decode response body: %w", err))
		return m
	}
	if len(body.Result) == 0 {
		m.Errors = append(m.Errors, fmt.Errorf("no transactions found"))
		return m
	}
	m.SuccessChecks++
	m.IndexingLatency = float64(time.Now().Unix() - body.Result[0].Utime)
	return m
}
