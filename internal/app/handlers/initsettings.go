package handlers

import (
	"github.com/kindenko/go-shorterurl/config"
	"github.com/kindenko/go-shorterurl/internal/app/storage"
)

type Handlers struct {
	cfg     *config.AppConfig
	storage storage.MyStorage
}

// часть мейна добить

func NewHandlers(cfg *config.AppConfig) *Handlers {

	c := cfg
	storage := storage.Init(c)

	return &Handlers{
		cfg:     cfg,
		storage: storage,
	}
}
