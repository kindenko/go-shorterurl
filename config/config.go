package config

import (
	"flag"
	"os"
	"strings"
)

const (
	DBhost     = "localhost"
	DBuser     = "postgres"
	DBpassword = ""
	DBdbname   = "postgres"
)

type AppConfig struct {
	Host           string `env:"SERVER_ADDRESS"`
	ResultURL      string `env:"BASE_URL"`
	FilePATH       string `env:"FILE_STORAGE_PATH"`
	DataBaseString string `env:"DATABASE_DSN"`
	// возможно добавить configDB
}

var cfg AppConfig

func NewCfg() *AppConfig {

	cfq := &AppConfig{}

	flag.StringVar(&cfq.Host, "a", "localhost:8080", "Host")
	flag.StringVar(&cfq.ResultURL, "b", "http://localhost:8080", "Result URL")
	flag.StringVar(&cfq.FilePATH, "f", "/tmp/short-url-db.json", "FilePATH")
	flag.StringVar(&cfq.DataBaseString, "d", "", "Connect to DB")

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
	if dbstring := os.Getenv("DATABASE_DSN"); dbstring != "" {
		cfq.DataBaseString = strings.TrimSpace(dbstring)
	}

	return cfq
}
