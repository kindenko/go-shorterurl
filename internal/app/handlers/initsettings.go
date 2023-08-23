package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kindenko/go-shorterurl/config"
	"github.com/kindenko/go-shorterurl/internal/app/logger"
	"github.com/kindenko/go-shorterurl/internal/app/storage"
	"github.com/kindenko/go-shorterurl/internal/app/zip"
)

type Handlers struct {
	cfg     *config.AppConfig
	storage storage.MyStorage
	mux     *chi.Mux
}

func NewHandlers(cfg *config.AppConfig, mux *chi.Mux) *Handlers {

	c := cfg
	storage := storage.Init(c)

	return &Handlers{
		cfg:     cfg,
		storage: storage,
		mux:     mux,
	}
}

func (h *Handlers) Init() {
	r := chi.NewRouter()
	r.Use(logger.WithLogging)
	r.Use(zip.MiddlewareCompressGzip)
	r.Use(middleware.Compress(5, "text/html",
		"application/x-gzip",
		"text/plain",
		"application/json"))

	r.Route("/", func(r chi.Router) {
		r.Post("/", h.PostHandler)
		r.Post("/api/shorten", h.PostJSONHandler)
		r.Get("/{shortUrl}", h.GetHandler)
		r.Get("/ping", h.Ping)
		r.Post("/api/shorten/batch", h.Batch)
	})
}

func (h *Handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
