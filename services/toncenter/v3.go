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
	token  string
}

func NewV3Monitoring(name, prefix, token string) *V3Monitoring {
	return &V3Monitoring{
		name:   name,
		prefix: prefix,
		token:  token,
	}
}

func (vm *V3Monitoring) GetMetrics(ctx context.Context) services.ApiMetrics {
	m := services.ApiMetrics{
		ServiceName: vm.name,
	}
	m.TotalChecks++
	t := time.Now()
	url := fmt.Sprintf("%v/account?address=%v", vm.prefix, config.ElectorAccountID.ToHuman(true, false))
	if vm.token != "" {
		url += fmt.Sprintf("&api_key=%v", vm.token)
	}
	r, err := http.Get(url)
	if err != nil {
		m.Errors = append(m.Errors, fmt.Errorf("failed to get account state: %w", err))
	} else if r.StatusCode != http.StatusOK {
		m.Errors = append(m.Errors, fmt.Errorf("invalid status code %v", r.StatusCode))
	} else {
		m.SuccessChecks++
	}
	m.HttpsLatency = time.Since(t).Seconds()
	m.TotalChecks++

	url = fmt.Sprintf("%v/transactions?account=%v&limit=1", vm.prefix, config.ElectorAccountID.ToHuman(true, false))
	if vm.token != "" {
		url += fmt.Sprintf("&api_key=%v", vm.token)
	}
	r, err = http.Get(url)
	if err != nil || r.StatusCode != http.StatusOK {
		m.Errors = append(m.Errors, fmt.Errorf("failed to get account transactions: %w", err))
		return m
	} else if r.StatusCode != http.StatusOK {
		m.Errors = append(m.Errors, fmt.Errorf("invalid status code %v", r.StatusCode))
		return m
	}
	defer r.Body.Close()
	var result struct {
		Transactions []struct {
			Now int64 `json:"now"`
		}
	}
	if err = json.NewDecoder(r.Body).Decode(&result); err != nil {
		m.Errors = append(m.Errors, fmt.Errorf("failed to decode response body: %w", err))
		return m
	}
	if len(result.Transactions) == 0 {
		m.Errors = append(m.Errors, fmt.Errorf("no transactions found"))
		return m
	}
	m.SuccessChecks++
	m.IndexingLatency = float64(time.Now().Unix() - result.Transactions[0].Now)
	return m
}
