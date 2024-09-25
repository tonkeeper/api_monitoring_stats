package tonxapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"api_monitoring_stats/config"
	"api_monitoring_stats/services"
)

type TonXAPIMonitoring struct {
	name string
	url  string
}

func NewTonXAPIMonitoring(name, url string) *TonXAPIMonitoring {
	return &TonXAPIMonitoring{
		name: name,
		url:  url,
	}
}

func (tm *TonXAPIMonitoring) GetMetrics(ctx context.Context) services.ApiMetrics {
	m := services.ApiMetrics{
		ServiceName: tm.name,
	}
	m.TotalChecks++
	t := time.Now()
	elector := config.ElectorAccountID.ToHuman(true, false)

  payloadData := map[string]interface{}{
        "id":      "1",
        "jsonrpc": "2.0",
        "method":  "getTransactions",
        "params": map[string]string{
            "account": elector,
			"sort": "DESC",
        },
    }


    payload, _ := json.Marshal(payloadData)
    fullURL := fmt.Sprintf("%s/%s", tm.url, config.Config.TonXAPIToken)    
    req, err := http.NewRequestWithContext(ctx, "POST", fullURL, bytes.NewReader(payload))
	if err != nil {
		m.Errors = append(m.Errors, fmt.Errorf("failed to create request: %w", err))
		return m
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	r, err := client.Do(req)
	if err != nil {
		m.Errors = append(m.Errors, fmt.Errorf("failed to get transactions: %w", err))
		return m
	}
	defer r.Body.Close()

	m.HttpsLatency = time.Since(t).Seconds()

	if r.StatusCode != http.StatusOK {
		m.Errors = append(m.Errors, fmt.Errorf("invalid status code %v", r.StatusCode))
		return m
	}

	var response struct {
		Result []struct {
			Utime int64 `json:"now"`
		} `json:"result"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err = json.NewDecoder(r.Body).Decode(&response); err != nil {
		m.Errors = append(m.Errors, fmt.Errorf("failed to decode response body: %w", err))
		return m
	}

	if response.Error != nil {
		m.Errors = append(m.Errors, fmt.Errorf("API error: code %d, message: %s", response.Error.Code, response.Error.Message))
		return m
	}

	if len(response.Result) == 0 {
		m.Errors = append(m.Errors, fmt.Errorf("no transactions found"))
		return m
	}

	m.SuccessChecks++
	m.IndexingLatency = float64(time.Now().Unix() - response.Result[0].Utime)
	return m
}