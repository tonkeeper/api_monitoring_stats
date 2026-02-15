package txpropagation

import (
	"context"
	"time"

	"github.com/tonkeeper/tongo/ton"
)

type LiteserverChecker struct {
	ServiceName string
	GetState    func(ctx context.Context, account ton.AccountID) (lastLT uint64, err error)
}

func (c *LiteserverChecker) Name() string { return c.ServiceName }

func (c *LiteserverChecker) WaitForTx(ctx context.Context, account ton.AccountID, initialLT uint64, sendTime time.Time) (float64, error) {
	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}
		lastLT, err := c.GetState(ctx, account)
		if err != nil {
			time.Sleep(50 * time.Millisecond)
			continue
		}
		if lastLT > initialLT {
			return time.Since(sendTime).Seconds(), nil
		}
		time.Sleep(50 * time.Millisecond)
	}
}
