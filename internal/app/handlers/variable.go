package handlers

import (
	"github.com/kindenko/go-shorterurl/config"
)

type Handlers struct {
	cfg *config.AppConfig
}

func NewHandlers(cfg *config.AppConfig) *Handlers {

	return &Handlers{
		cfg: cfg,
	}
}
