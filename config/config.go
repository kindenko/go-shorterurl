package config

import (
	"flag"
	"os"
	"strings"
)

type AppConfig struct {
	Host      string
	ResultURL string
}

func NewCfg() *AppConfig {

	cfq := &AppConfig{}

	flag.StringVar(&cfq.Host, "a", "localhost:8080", "It's a Host")
	flag.StringVar(&cfq.ResultURL, "b", "http://localhost:8080", "It's a Result URL")

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		cfq.ResultURL = strings.TrimSpace(baseURL)
	}
	if host := os.Getenv("SERVER_ADDRESS"); host != "" {
		cfq.Host = strings.TrimSpace(host)
	}

	return cfq
}
