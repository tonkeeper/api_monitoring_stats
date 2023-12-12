package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/tonkeeper/opentonapi/tonapi"
	"github.com/tonkeeper/tongo/ton"
)

var electorAccountID = ton.MustParseAccountID("Ef8zMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzM0vF")

type service string

const (
	tonApi    service = "TonAPI"
	dton      service = "dTON"
	tonCenter service = "TonCenter"
)

var (
	metricServiceTimeHistogramVec = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "service_functions_time",
			Help:    "Service functions execution duration distribution in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 1, 5, 10, 60},
		},
		[]string{"service", "method"},
	)

	metricServiceRequestSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "service_request_success",
		Help: "The total number of success request",
	}, []string{"service", "method"})

	metricServiceRequestFails = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "service_request_fails",
		Help: "The total number of failed request",
	}, []string{"service", "method"})
)

type Metric struct {
	tonApiClient *tonapi.Client
}

func NewMetric(tonApiClient *tonapi.Client) *Metric {
	return &Metric{
		tonApiClient: tonApiClient,
	}
}

type metric interface {
	GetMetrics(service service)
}

func (m *Metric) GetMetrics(service service) {
	switch service {
	case tonApi:
		m.tonApiMetric()
	}
}

func (m *Metric) tonApiMetric() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	accountState := func() {
		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			metricServiceTimeHistogramVec.WithLabelValues(string(tonApi), "account_state").Observe(v)
		}))
		defer timer.ObserveDuration()

		if _, err := m.tonApiClient.GetAccountState(ctx, electorAccountID); err == nil {
			metricServiceRequestSuccess.WithLabelValues(string(tonApi), "account_state").Inc()
		} else {
			metricServiceRequestFails.WithLabelValues(string(tonApi), "account_state").Inc()
		}
	}

	accountTransactions := func() {
		now := time.Now().Unix()

		transactions, err := m.tonApiClient.GetBlockchainAccountTransactions(ctx, tonapi.GetBlockchainAccountTransactionsParams{
			AccountID: electorAccountID.ToRaw(),
			Limit:     tonapi.NewOptInt32(100),
		})
		if err != nil {
			metricServiceRequestFails.WithLabelValues(string(tonApi), "account_transactions").Inc()
		} else {
			metricServiceRequestSuccess.WithLabelValues(string(tonApi), "account_transactions").Inc()
		}
		if len(transactions.Transactions) == 0 {
			return
		}

		timeOfLastTransaction := now - transactions.Transactions[0].Utime
		fmt.Println(timeOfLastTransaction) // TODO: add to metric
	}

	accountState()
	accountTransactions()
}
