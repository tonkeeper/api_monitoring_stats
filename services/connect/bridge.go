package connect

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"api_monitoring_stats/services"
)

type Bridge struct {
	name             string
	url              string
	id               string
	connected        bool
	reconnectCounter int
	data             chan string
}

func randString() string {
	var b [32]byte
	rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func NewBridge(name, url string) *Bridge {
	b := &Bridge{
		name: name,
		url:  url,
		data: make(chan string, 10),
		id:   randString(),
	}
	go b.connect()
	return b
}

func (b *Bridge) GetMetrics(ctx context.Context) services.BridgeMetrics {
	m := services.BridgeMetrics{
		ServiceName: b.name,
		TotalChecks: 1,
		Reconnects:  b.reconnectCounter,
	}
	if !b.connected {
		return m
	}
	t := time.Now()
	payload := randString()
	resp, err := http.Post(fmt.Sprintf("%s/message?client_id=%s&to=%s&ttl=300", b.url, randString(), b.id), "text/plain", strings.NewReader(payload))
	if err != nil || resp.StatusCode != 200 {
		return m
	}

	timer := time.NewTimer(time.Second * 10)
external:
	for {
		select {
		case <-timer.C:
			break external
		case data := <-b.data:
			var message struct {
				Message string
			}
			err := json.Unmarshal([]byte(data), &message)
			if err != nil {
				m.Errors = append(m.Errors, err)
				return m
			}
			if message.Message == payload {
				m.SuccessChecks += 1
				break external
			}
		}

	}
	m.TransferLatency = time.Since(t).Seconds()
	return m
}

func (b *Bridge) connect() {
	for {
		resp, err := http.Get(b.url + "/events?client_id=" + b.id)
		if err != nil || resp.StatusCode != 200 {
			b.connected = false
			time.Sleep(time.Second * 10)
			continue
		} else {
			b.connected = true
		}
		for {
			var event string
			var data string
			_, err := fmt.Fscanf(resp.Body, "%s %s", &event, &data)
			if err != nil {
				if err.Error() == "unexpected newline" {
					continue
				}
				fmt.Println("bridge", b.name, err)
				b.connected = false
				b.reconnectCounter++
				break
			}
			if event == "data:" {
				select {
				case b.data <- data:
				default:

				}
			}
		}
		time.Sleep(time.Second * 10)
	}
}
