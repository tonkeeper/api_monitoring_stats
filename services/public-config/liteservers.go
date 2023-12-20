package public_config

import (
	"api_monitoring_stats/config"
	"api_monitoring_stats/services"
	"context"
	"fmt"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"time"
)

type LiteServersMetrics struct {
	client *ton.APIClient
}

func NewLiteServersMetrics() *LiteServersMetrics {
	l := &LiteServersMetrics{}
	go func() {
		for l.client == nil {
			err := l.connect()
			if err != nil {
				time.Sleep(time.Second * 10)
			}
		}
	}()
	return l
}

func (lm *LiteServersMetrics) connect() error {
	pool := liteclient.NewConnectionPool()
	configUrl := "https://api.tontech.io/ton/wallet-mainnet.autoconf.json"
	err := pool.AddConnectionsFromConfigUrl(context.Background(), configUrl)
	if err != nil {
		return err
	}
	lm.client = ton.NewAPIClient(pool)

	return nil

}

func (lm *LiteServersMetrics) GetMetrics(ctx context.Context) services.ApiMetrics {
	m := services.ApiMetrics{
		ServiceName: "public liteservers",
	}

	m.TotalChecks++
	if lm.client == nil {
		m.Errors = append(m.Errors, liteclient.ErrNoConnections)
		return m
	}

	t := time.Now()

	b, err := lm.client.GetMasterchainInfo(ctx)
	if err != nil {
		m.Errors = append(m.Errors, err)
	}

	elector := address.MustParseAddr(config.ElectorAccountID.ToHuman(true, false))
	a, err := lm.client.GetAccount(ctx, b, elector)
	if err != nil {
		m.Errors = append(m.Errors, err)
		return m
	} else if a.State == nil || a.State.Balance.Nano().Int64() == 0 {
		m.Errors = append(m.Errors, fmt.Errorf("invalid account state"))
	} else {
		m.SuccessChecks++
	}
	m.HttpsLatency = time.Since(t).Seconds()

	m.TotalChecks++
	txs, err := lm.client.ListTransactions(ctx, elector, 1, a.LastTxLT, a.LastTxHash)
	if err != nil {
		m.Errors = append(m.Errors, err)
		return m
	}
	if len(txs) == 0 {
		m.Errors = append(m.Errors, fmt.Errorf("invalid txs"))
		return m
	}
	m.IndexingLatency = float64(t.Unix()) - float64(txs[0].Now)
	m.SuccessChecks++
	return m
}
