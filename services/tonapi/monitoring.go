package tonapi

import (
	"context"
	"log"
	"time"

	"api_monitoring_stats/config"
	"api_monitoring_stats/services"
	"github.com/tonkeeper/opentonapi/tonapi"
)

type TonAPI struct{}

func NewMonitoring() *TonAPI {
	return &TonAPI{}
}

func (t *TonAPI) GetMetrics() (services.ApiMetrics, error) {
	tonApiClient, err := tonapi.New()
	if err != nil {
		log.Printf("failed to init tonapi client: %v", err)
		return services.ApiMetrics{}, err
	}

	metrics := services.ApiMetrics{
		ServiceName: services.TonApi,
		Alive:       true,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start := time.Now()
	_, err = tonApiClient.GetAccountState(ctx, config.ElectorAccountID)
	if err != nil {
		log.Printf("failed to get account state, service: %v, %v", services.TonApi, err)
		metrics.Alive = false
	}
	metrics.HttpsLatency = time.Since(start).Seconds()

	transactions, err := tonApiClient.GetBlockchainAccountTransactions(ctx, tonapi.GetBlockchainAccountTransactionsParams{
		AccountID: config.ElectorAccountID.ToRaw(),
		Limit:     tonapi.NewOptInt32(10),
	})
	if err != nil {
		log.Printf("failed to get account transactions, service: %v, %v", services.TonApi, err)
		metrics.Alive = false
		return metrics, err
	}
	if len(transactions.Transactions) == 0 {
		return metrics, nil
	}
	metrics.IndexingLatency = float64(time.Now().Unix() - transactions.Transactions[0].Utime)

	return metrics, nil
}
