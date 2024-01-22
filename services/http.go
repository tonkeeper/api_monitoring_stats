package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func HttpGet(ctx context.Context, totalChecks, successChecks *int, errors *[]error, url string, respObject any) float64 {
	*totalChecks += 1
	t := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		*errors = append(*errors, fmt.Errorf("can't create request to %v: %w", url, err))
		return 0
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		*errors = append(*errors, fmt.Errorf("can't get %v: %w", url, err))
		return 0
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		*errors = append(*errors, fmt.Errorf("invalid status code: %d for %v", resp.StatusCode, url))
		return 0
	}
	if respObject != nil {
		err = json.NewDecoder(resp.Body).Decode(respObject)
		if err != nil {
			*errors = append(*errors, fmt.Errorf("can't process response from %v: %w", url, err))
			return 0
		}
	}
	*successChecks += 1
	return time.Since(t).Seconds()
}

func HttpPost(ctx context.Context, totalChecks, successChecks *int, errors *[]error, url string, body io.Reader, respObject any) float64 {
	*totalChecks += 1
	t := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		*errors = append(*errors, fmt.Errorf("can't create request to %v: %w", url, err))
		return 0
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		*errors = append(*errors, fmt.Errorf("can't post %v: %w", url, err))
		return 0
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		*errors = append(*errors, fmt.Errorf("invalid status code: %d for %v", resp.StatusCode, url))
		return 0
	}
	if respObject != nil {
		err = json.NewDecoder(resp.Body).Decode(respObject)
		if err != nil {
			*errors = append(*errors, fmt.Errorf("can't process response from %v: %w", url, err))
			return 0
		}
	}
	*successChecks += 1
	return time.Since(t).Seconds()
}
