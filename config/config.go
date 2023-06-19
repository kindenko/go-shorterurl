package config

import (
	"flag"
	"os"
	"strings"
)

var Config struct {
	Host      string
	ResultURL string
}

func init() {

	if host := os.Getenv("SERVER_ADDRESS"); host != "" {
		Config.Host = strings.TrimSpace(host)
	}
	flag.StringVar(&Config.Host, "a", "localhost:8080", "It's a Host")

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		Config.ResultURL = strings.TrimSpace(baseURL)
	}
	flag.StringVar(&Config.ResultURL, "b", "http://localhost:8080", "It's a Result URL")
}
