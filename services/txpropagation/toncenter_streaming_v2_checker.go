package txpropagation

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tonkeeper/tongo/ton"
)

type ToncenterStreamingV2Checker struct {
	ServiceName string
	BaseURL     string
	Token       string
	MinFinality string
}

func (c *ToncenterStreamingV2Checker) Name() string { return c.ServiceName }

func (c *ToncenterStreamingV2Checker) WaitForTx(ctx context.Context, account ton.AccountID, initialLT uint64, sendTime time.Time) (float64, error) {
	body := map[string]interface{}{
		"addresses":            []string{account.ToRaw()},
		"types":                []string{"transactions"},
		"min_finality":         c.MinFinality,
		"include_address_book": false,
		"include_metadata":     false,
	}
	bodyJSON, _ := json.Marshal(body)

	var resp *http.Response
	for retry := 0; retry < 3; retry++ {
		if retry > 0 {
			backoff := time.Duration(retry*10) * time.Second
			select {
			case <-ctx.Done():
				return 0, ctx.Err()
			case <-time.After(backoff):
			}
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/sse", bytes.NewReader(bodyJSON))
		if err != nil {
			return 0, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "text/event-stream")
		if c.Token != "" {
			req.Header.Set("X-Api-Key", c.Token)
		}

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return 0, err
		}
		if resp.StatusCode == http.StatusOK {
			break
		}
		if resp.StatusCode == 429 {
			resp.Body.Close()
			if retry == 2 {
				return 0, fmt.Errorf("streaming sse status 429 (rate limited) after retries")
			}
			continue
		}
		resp.Body.Close()
		return 0, fmt.Errorf("streaming sse status %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	const maxLineSize = 1024 * 1024
	scanner.Buffer(nil, maxLineSize)

	var dataLines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			dataLines = append(dataLines, strings.TrimPrefix(line, "data: "))
			continue
		}

		if len(dataLines) > 0 {
			data := strings.Join(dataLines, "\n")
			dataLines = nil
			var msg struct {
				Type         string `json:"type"`
				Finality     string `json:"finality"`
				Transactions []struct {
					Account string      `json:"account"`
					Lt      interface{} `json:"lt"`
				} `json:"transactions"`
			}
			if err := json.Unmarshal([]byte(data), &msg); err != nil {
				select {
				case <-ctx.Done():
					return 0, ctx.Err()
				default:
				}
				continue
			}
			if msg.Type != "transactions" || msg.Finality != c.MinFinality {
				select {
				case <-ctx.Done():
					return 0, ctx.Err()
				default:
				}
				continue
			}
			for _, tx := range msg.Transactions {
				lt, err := parseLT(tx.Lt)
				if err != nil {
					continue
				}
				if lt > initialLT {
					return time.Since(sendTime).Seconds(), nil
				}
			}
		}
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return 0, ctx.Err()
}

func parseLT(v interface{}) (uint64, error) {
	switch x := v.(type) {
	case string:
		return strconv.ParseUint(x, 10, 64)
	case float64:
		if x < 0 {
			return 0, strconv.ErrSyntax
		}
		return uint64(x), nil
	case nil:
		return 0, strconv.ErrSyntax
	default:
		return 0, strconv.ErrSyntax
	}
}
