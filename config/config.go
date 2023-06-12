package config

import (
	"flag"
	"os"
	"strings"
)

var SetConfig struct {
	Host      string
	ResultURL string
}

func init() {

	if host := os.Getenv("SERVER_ADDRESS"); host != "" {
		SetConfig.Host = strings.TrimSpace(host)
	}
	flag.StringVar(&SetConfig.Host, "a", "localhost:8080", "It's a Host")

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		SetConfig.ResultURL = strings.TrimSpace(baseURL)
	}
	flag.StringVar(&SetConfig.ResultURL, "b", "http://localhost:8080", "It's a Result URL")
}
