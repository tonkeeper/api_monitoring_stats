package tonapi

import (
	"api_monitoring_stats/services"
	"context"
	"fmt"
	"time"

	"api_monitoring_stats/config"
	"github.com/tonkeeper/tonapi-go"
)

type TonAPI struct{}

const serviceName = "tonapi"

func NewMonitoring() *TonAPI {
	return &TonAPI{}
}

var tonApiClient *tonapi.Client

func init() {
	var err error
	tonApiClient, err = tonapi.New()
	if err != nil {
		panic(err)
	}
}

func (t *TonAPI) GetMetrics(ctx context.Context) services.ApiMetrics {

	metrics := services.ApiMetrics{
		ServiceName: serviceName,
	}
	start := time.Now()
	metrics.TotalChecks++
	_, err := tonApiClient.GetAccountState(ctx, config.ElectorAccountID)
	if err != nil {
		metrics.Errors = append(metrics.Errors, fmt.Errorf("failed to get account state: %w", err))
	} else {
		metrics.SuccessChecks++
	}
	metrics.HttpsLatency = time.Since(start).Seconds()

	metrics.TotalChecks++
	transactions, err := tonApiClient.GetBlockchainAccountTransactions(ctx, tonapi.GetBlockchainAccountTransactionsParams{
		AccountID: config.ElectorAccountID.ToRaw(),
		Limit:     tonapi.NewOptInt32(1),
	})
	if err != nil {
		metrics.Errors = append(metrics.Errors, fmt.Errorf("failed to get account transactions: %w", err))
		return metrics
	}
	if len(transactions.Transactions) == 0 {
		metrics.Errors = append(metrics.Errors, fmt.Errorf("no transactions found"))
		return metrics
	}
	metrics.SuccessChecks++
	metrics.IndexingLatency = float64(time.Now().Unix() - transactions.Transactions[0].Utime)

	return metrics
}
