package txpropagation

import (
	"context"
	"time"

	"github.com/tonkeeper/tongo/ton"
)

type Checker interface {
	Name() string
	WaitForTx(ctx context.Context, account ton.AccountID, initialLT uint64, sendTime time.Time) (latencySec float64, err error)
}
