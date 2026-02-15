package txpropagation

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"api_monitoring_stats/services"

	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/wallet"
)

type Runner struct {
	Seed       string
	UseTestnet bool
	Checkers   []Checker
	Interval   time.Duration
}

func (r *Runner) GetMetrics(ctx context.Context) services.TxPropagationMetrics {
	latencies, err := r.Run(ctx)
	if err != nil {
		return services.TxPropagationMetrics{Errors: []error{err}}
	}
	return services.TxPropagationMetrics{ServiceName: "tx-propagation", Latency: latencies["liteserver"], Errors: []error{err}}
}

func (r *Runner) CheckInterval() time.Duration {
	if r.Interval > 0 {
		return r.Interval
	}
	return 5 * time.Minute
}

func (r *Runner) Timeout() time.Duration {
	return 1 * time.Minute
}

// Run sends one tx and returns latency per service (seconds). Checkers are run in parallel.
func (r *Runner) Run(ctx context.Context) (latencies map[string]float64, err error) {
	latencies = make(map[string]float64)

	var opts liteapi.Option = liteapi.Mainnet()
	if r.UseTestnet {
		opts = liteapi.Testnet()
	}
	cli, err := liteapi.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("liteapi client: %w", err)
	}

	privateKey, err := wallet.SeedToPrivateKey(r.Seed)
	if err != nil {
		return nil, fmt.Errorf("seed: %w", err)
	}

	w, err := wallet.New(privateKey, wallet.V3R2, 0, nil, cli)
	if err != nil {
		return nil, fmt.Errorf("wallet: %w", err)
	}

	walletAddress := w.GetAddress()
	state, err := cli.GetAccountState(ctx, walletAddress)
	if err != nil {
		return nil, fmt.Errorf("get account state: %w", err)
	}
	initialLT := state.LastTransLt

	amount := uint64(10_000_000) // 0.01 TON
	msg := wallet.SimpleTransfer{
		Amount:  tlb.Grams(amount),
		Address: walletAddress, // will be bounced back to the sender
		Comment: fmt.Sprintf("tx-timing %d", time.Now().UTC().Unix()),
	}

	sendTime := time.Now()
	if err = w.Send(ctx, msg); err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}
	log.Printf("[tx-propagation] tx sent, initialLT=%d", initialLT)

	allCheckers := append([]Checker{}, r.Checkers...)
	allCheckers = append(allCheckers, &LiteserverChecker{
		ServiceName: "liteserver",
		GetState: func(ctx context.Context, account ton.AccountID) (uint64, error) {
			st, err := cli.GetAccountState(ctx, account)
			if err != nil {
				return 0, err
			}
			return st.LastTransLt, nil
		},
	})

	ctxWait, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, ch := range allCheckers {
		ch := ch
		wg.Add(1)
		go func() {
			defer wg.Done()
			lat, err := ch.WaitForTx(ctxWait, walletAddress, initialLT, sendTime)
			mu.Lock()
			if err != nil {
				log.Printf("[tx-propagation] %s: %v", ch.Name(), err)
			} else {
				latencies[ch.Name()] = lat
			}
			mu.Unlock()
		}()
	}
	wg.Wait()

	return latencies, nil
}
