package txpropagation

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/tonkeeper/tongo/ton"
)

type TonapiChecker struct {
	ServiceName string
	BaseURL     string
	APIKey      string
}

func (c *TonapiChecker) Name() string { return c.ServiceName }

func (c *TonapiChecker) WaitForTx(ctx context.Context, account ton.AccountID, initialLT uint64, sendTime time.Time) (float64, error) {
	accountRaw := account.ToRaw()
	url := c.BaseURL + "/v2/blockchain/accounts/" + accountRaw + "/transactions?limit=1"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}
	client := &http.Client{Timeout: 15 * time.Second}
	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}
		resp, err := client.Do(req)
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		var body struct {
			Transactions []struct {
				Lt int64 `json:"lt"`
			} `json:"transactions"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&body)
		resp.Body.Close()
		if len(body.Transactions) > 0 && uint64(body.Transactions[0].Lt) > initialLT {
			return time.Since(sendTime).Seconds(), nil
		}
		time.Sleep(200 * time.Millisecond)
	}
}
