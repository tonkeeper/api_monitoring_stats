package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func HttpGet(totalChecks, successChecks *int, errors *[]error, url string, respObject any) float64 {
	*totalChecks += 1
	t := time.Now()
	resp, err := http.Get(url)
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

func HttpPost(totalChecks, successChecks *int, errors *[]error, url string, body io.Reader, respObject any) float64 {
	*totalChecks += 1
	t := time.Now()
	resp, err := http.Post(url, "application/json", body)
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
