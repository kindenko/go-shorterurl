package config

import (
	"flag"
	"os"
	"strings"
)

type AppConfig struct {
	Host      string `env:"SERVER_ADDRESS"`
	ResultURL string `env:"BASE_URL"`
	FilePATH  string `env:"FILE_STORAGE_PATH"`
}

func NewCfg() *AppConfig {

	cfq := &AppConfig{}

	flag.StringVar(&cfq.Host, "a", "localhost:8080", "It's a Host")
	flag.StringVar(&cfq.ResultURL, "b", "http://localhost:8080", "It's a Result URL")
	flag.StringVar(&cfq.FilePATH, "f", "/tmp/short-url-db.json", "It's a FilePATH")

	flag.Parse()

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		cfq.ResultURL = strings.TrimSpace(baseURL)
	}
	if host := os.Getenv("SERVER_ADDRESS"); host != "" {
		cfq.Host = strings.TrimSpace(host)
	}
	if filepath := os.Getenv("FILE_STORAGE_PATH"); filepath != "" {
		cfq.FilePATH = strings.TrimSpace(filepath)
	}

	return cfq
}
