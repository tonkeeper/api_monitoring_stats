package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

var Config = struct {
	Dev         string `env:"DEV"`
	Port        int    `env:"PORT" envDefault:"8888"`
	MetricsPort int    `env:"METRICS_PORT" envDefault:"9010"`
	TonApiToken string `env:"TONAPI_TOKEN"`
}{}

func LoadConfig() {
	if err := env.Parse(&Config); err != nil {
		log.Fatalf("❗️failed to parse config: %v\n", err)
	}
}
