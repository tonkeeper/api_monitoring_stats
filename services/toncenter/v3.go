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

type V3Monitoring struct {
	name   string
	prefix string
}

func NewV3Monitoring(name, prefix string) *V3Monitoring {
	return &V3Monitoring{
		name:   name,
		prefix: prefix,
	}
}

func (vm *V3Monitoring) GetMetrics(ctx context.Context) services.ApiMetrics {
	m := services.ApiMetrics{
		ServiceName: vm.name,
	}
	time.Sleep(1 * time.Second)
	m.TotalChecks++
	t := time.Now()
	r, err := http.Get(fmt.Sprintf("%v/account?address=%v", vm.prefix, config.ElectorAccountID.ToHuman(true, false)))
	if err != nil || r.StatusCode != http.StatusOK {
		m.Errors = append(m.Errors, fmt.Errorf("failed to get account state: %w, status code: %v", err, r.StatusCode))
	} else {
		m.SuccessChecks++
	}
	m.HttpsLatency = time.Since(t).Seconds()
	m.TotalChecks++

	time.Sleep(1 * time.Second)
	r, err = http.Get(fmt.Sprintf("%v/transactions?account=%v&limit=1", vm.prefix, config.ElectorAccountID.ToHuman(true, false)))
	if err != nil {
		m.Errors = append(m.Errors, fmt.Errorf("failed to get account transactions: %w", err))
		return m
	}
	defer r.Body.Close()

	var result []struct {
		Now int64 `json:"now"`
	}
	if err = json.NewDecoder(r.Body).Decode(&result); err != nil {
		m.Errors = append(m.Errors, fmt.Errorf("failed to decode response body: %w", err))
		return m
	}
	if len(result) == 0 {
		m.Errors = append(m.Errors, fmt.Errorf("no transactions found"))
		return m
	}
	m.SuccessChecks++
	m.IndexingLatency = float64(time.Now().Unix() - result[0].Now)
	return m
}
