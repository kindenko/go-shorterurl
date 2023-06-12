package config

import (
	"flag"
	"os"
)

var SetConfig struct {
	Host      string
	ResultURL string
}

func init() {

	if host := os.Getenv("SERVER_ADDRESS"); host != "" {
		SetConfig.Host = host
	} else {
		flag.StringVar(&SetConfig.Host, "a", "localhost:8080", "It's a Host")
	}

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		SetConfig.ResultURL = baseURL
	} else {
		flag.StringVar(&SetConfig.ResultURL, "b", "http://localhost:8080", "It's a Result URL")
	}

}
