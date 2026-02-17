package txpropagation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tonkeeper/tongo/ton"
)

type ToncenterChecker struct {
	ServiceName string
	Prefix      string
	Token       string
}

func (c *ToncenterChecker) Name() string { return c.ServiceName }

func (c *ToncenterChecker) WaitForTx(ctx context.Context, account ton.AccountID, initialLT uint64, sendTime time.Time) (float64, error) {
	url := fmt.Sprintf("%s/transactions?account=%s&limit=3&offset=0&sort=desc", c.Prefix, account.ToHuman(true, false))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	if c.Token != "" {
		req.Header.Add("X-Api-Key", c.Token)
	}
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-ticker.C:
			reqClone := req.Clone(ctx)
			resp, err := http.DefaultClient.Do(reqClone)
			if err != nil {
				continue
			}
			if resp.StatusCode != http.StatusOK {
				resp.Body.Close()
				continue
			}
			var body struct {
				Transactions []struct {
					Lt uint64 `json:"lt,string"`
				} `json:"transactions"`
			}
			err = json.NewDecoder(resp.Body).Decode(&body)
			resp.Body.Close()
			if err != nil {
				if ctx.Err() != nil {
					return 0, ctx.Err()
				}
				continue
			}
			if len(body.Transactions) == 0 {
				continue
			}
			currentLT := body.Transactions[0].Lt
			if currentLT > initialLT {
				return time.Since(sendTime).Seconds(), nil
			}
		}
	}
}
