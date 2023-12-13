package dton

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"api_monitoring_stats/config"
	"api_monitoring_stats/services"
)

type Monitoring struct {
	name   string
	prefix string
}

type graphQLQuery struct {
	Query string `json:"query"`
}

func NewMonitoring(name, prefix string) *Monitoring {
	return &Monitoring{
		name:   name,
		prefix: prefix,
	}
}

func (m *Monitoring) GetMetrics(ctx context.Context) services.ApiMetrics {
	metrics := services.ApiMetrics{
		ServiceName: m.prefix,
	}
	start := time.Now()
	metrics.TotalChecks++
	query := fmt.Sprintf(`{account_states(account: {address_friendly: "%v"}) {account_storage_balance_grams}}`, config.ElectorAccountID.ToHuman(true, false))
	body, _ := json.Marshal(graphQLQuery{Query: query})
	if _, err := sendRequest(m.prefix, body); err != nil {
		metrics.Errors = append(metrics.Errors, fmt.Errorf("failed to get account state: %w", err))
	} else {
		metrics.SuccessChecks++
	}
	metrics.HttpsLatency = time.Since(start).Seconds()

	metrics.TotalChecks++
	query = fmt.Sprintf(`{transactions(order_desc: true, page_size: 1, address_friendly: "%v"){gen_utime}}`, config.ElectorAccountID.ToHuman(true, false))
	body, _ = json.Marshal(graphQLQuery{Query: query})
	responseBody, err := sendRequest(m.prefix, body)
	if err != nil {
		metrics.Errors = append(metrics.Errors, fmt.Errorf("failed to get account transactions: %w", err))
		return metrics
	}
	var result struct {
		Data struct {
			Transactions []struct {
				GenUtime string `json:"gen_utime"`
			} `json:"transactions"`
		} `json:"data"`
	}
	if err = json.Unmarshal(responseBody, &result); err != nil {
		metrics.Errors = append(metrics.Errors, fmt.Errorf("failed to decode response body: %w", err))
		return metrics
	}
	if len(result.Data.Transactions) == 0 {
		metrics.Errors = append(metrics.Errors, fmt.Errorf("no transactions found"))
		return metrics
	}
	metrics.SuccessChecks++

	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	parsedTime, err := time.ParseInLocation("2006-01-02T15:04:05", result.Data.Transactions[0].GenUtime, moscowLocation)
	if err != nil {
		metrics.Errors = append(metrics.Errors, fmt.Errorf("failed to parse time"))
		return metrics
	}
	metrics.IndexingLatency = float64(time.Now().Unix() - parsedTime.Unix())
	return metrics
}

func sendRequest(url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("")
	}

	return respBody, nil
}
