package tonhub

import (
	"api_monitoring_stats/config"
	"api_monitoring_stats/services"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type V4Monitoring struct {
	prefix string
	name   string
}

func NewV4Monitoring(name, prefix string) *V4Monitoring {
	return &V4Monitoring{
		prefix: prefix,
		name:   name,
	}
}

func (vm *V4Monitoring) GetMetrics(ctx context.Context) services.ApiMetrics {
	m := services.ApiMetrics{
		ServiceName: vm.name,
	}

	m.TotalChecks++
	r, err := http.Get(vm.prefix + "/block/latest")
	if err != nil || r.StatusCode != http.StatusOK {
		m.Errors = append(m.Errors, fmt.Errorf("failed to get latest block: %w", err))
		return m
	}
	var block struct {
		Last struct {
			Seqno int64 `json:"seqno"`
		} `json:"last"`
		Now int64
	}
	err = json.NewDecoder(r.Body).Decode(&block)
	if err != nil {
		m.Errors = append(m.Errors, fmt.Errorf("failed to decode response body: %w", err))
		return m
	}
	r.Body.Close()
	t := time.Now()
	r, err = http.Get(fmt.Sprintf("%v/block/%v/%v", vm.prefix, block.Last.Seqno, config.ElectorAccountID.ToHuman(true, false)))
	if err != nil || r.StatusCode != http.StatusOK {
		m.Errors = append(m.Errors, fmt.Errorf("failed to get account state: %w", err))
	} else {
		m.SuccessChecks++
	}
	m.HttpsLatency = time.Since(t).Seconds()
	m.IndexingLatency = float64(time.Now().Unix() - block.Now)
	return m
}
