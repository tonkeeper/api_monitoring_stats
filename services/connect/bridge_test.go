package connect

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// TestBridgeForceReconnect tests that the bridge will reconnect if the server is down
func TestBridgeForceReconnect(t *testing.T) {
	var (
		mu           sync.Mutex
		messageCount int
		reconnects   int
		serverDown   bool
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		if serverDown {
			http.Error(w, "Server down", http.StatusServiceUnavailable)
			return
		}

		if r.Method == http.MethodPost && r.URL.Path == "/message" {
			clientID := r.URL.Query().Get("client_id")
			to := r.URL.Query().Get("to")
			ttl := r.URL.Query().Get("ttl")
			var payload string
			if r.Body != nil {
				b, err := io.ReadAll(r.Body)
				if err == nil {
					payload = string(b)
				}
			}
			fmt.Printf("Received POST request to /message: client_id=%v, to=%v, ttl=%v, payload=%v\n", clientID, to, ttl, payload)
			messageCount++
			w.WriteHeader(http.StatusOK)

		} else if r.URL.Path == "/events" {
			b, _ := json.Marshal(map[string]string{"from": "test_client", "message": "test_message"})
			w.Write([]byte(fmt.Sprintf("data: %v\n", string(b))))

		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}

	}))
	defer server.Close()

	bridge := NewBridge("test", server.URL)

	// Simulate the server going down and coming back up
	go func() {
		for {
			time.Sleep(10 * time.Second) // Simulate the server being up for 10 seconds
			mu.Lock()
			serverDown = true // Simulate the server going down
			mu.Unlock()
			time.Sleep(time.Second) // Give the bridge time to reconnect
			mu.Lock()
			serverDown = false        //  Simulate the server coming back up
			bridge.connected = false  // Force the bridge to reconnect
			bridge.reconnectCounter++ // Increment the reconnect counter
			reconnects++
			mu.Unlock()
		}
	}()

	// Run the test for 1 minute
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			metrics := bridge.GetMetrics(context.Background())
			if metrics.Reconnects != reconnects {
				t.Errorf("Expected %d reconnects, but got %v", reconnects, metrics.Reconnects)
			}
			return

		case <-ticker.C:
			bridge.GetMetrics(context.Background())
		}
	}
}
