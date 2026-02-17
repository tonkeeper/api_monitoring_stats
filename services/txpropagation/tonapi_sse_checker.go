package txpropagation

import (
	"context"
	"fmt"
	"time"

	"github.com/r3labs/sse/v2"
	"github.com/tonkeeper/tongo/ton"
)

type TonapiSSEChecker struct {
	ServiceName string
	BaseURL     string
	APIKey      string
}

func (c *TonapiSSEChecker) Name() string { return c.ServiceName }

func (c *TonapiSSEChecker) WaitForTx(ctx context.Context, account ton.AccountID, initialLT uint64, sendTime time.Time) (float64, error) {
	accountRaw := account.ToRaw()
	url := fmt.Sprintf("%s/sse/transactions?account=%s", c.BaseURL, accountRaw)

	foundCh := make(chan time.Time, 1)
	go listenSSE(ctx, url, c.APIKey, foundCh)

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-foundCh:
		return time.Since(sendTime).Seconds(), nil
	}
}

func listenSSE(ctx context.Context, url, apiKey string, foundCh chan<- time.Time) {
	client := sse.NewClient(url)
	if apiKey != "" {
		client.Headers = map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", apiKey),
		}
	}
	_ = client.SubscribeWithContext(ctx, "", func(msg *sse.Event) {
		switch string(msg.Event) {
		case "heartbeat":
			return
		case "message":
			select {
			case foundCh <- time.Now():
			default:
			}
		}
	})
}
