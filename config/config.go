package config

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/tonkeeper/tongo/ton"
	"github.com/xssnick/tonutils-go/liteclient"
)

var ElectorAccountID = ton.MustParseAccountID("Ef8zMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzM0vF")

var Config = struct {
	TonCenterApiToken string                        `env:"TONCENTER_API_TOKEN"`
	GetBlockKey       string                        `env:"GETBLOCK_KEY"`
	MetricsPort       int                           `env:"METRICS_PORT" envDefault:"9010"`
	DtonLiteServers   []liteclient.LiteserverConfig `env:"DTON_LITE_SERVERS"`
	DtonToken         string                        `env:"DTONTOKEN"`
	TonXAPIToken      string                        `env:"TONX_API_TOKEN"`
}{}

func LoadConfig() {
	err := env.ParseWithFuncs(&Config, map[reflect.Type]env.ParserFunc{
		reflect.TypeOf([]liteclient.LiteserverConfig{}): func(v string) (interface{}, error) {
			var servers []liteclient.LiteserverConfig
			for _, s := range strings.Split(v, ",") {
				if s == "" {
					continue
				}
				params := strings.Split(s, ":")
				if len(params) != 3 {
					return nil, fmt.Errorf("invalid liteserver string: %v", s)
				}
				ip := net.ParseIP(params[0])
				if ip == nil {
					return nil, fmt.Errorf("invalid lite server ip")
				}
				if ip.To4() == nil {
					return nil, fmt.Errorf("IPv6 not supported")
				}
				port, err := strconv.ParseInt(params[1], 10, 32)
				if err != nil {
					return nil, fmt.Errorf("invalid lite server port: %v", params[1])
				}
				servers = append(servers, liteclient.LiteserverConfig{
					IP:   int64(binary.BigEndian.Uint32(ip.To4())),
					Port: int(port),
					ID: liteclient.ServerID{
						Type: "pub.ed25519",
						Key:  params[2],
					},
				})
			}
			return servers, nil
		},
	})
	if err != nil {
		log.Fatalf("❗️failed to parse config: %v\n", err)
	}
}
