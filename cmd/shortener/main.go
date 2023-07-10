package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kindenko/go-shorterurl/config"
	"github.com/kindenko/go-shorterurl/internal/app/handlers"
	"github.com/kindenko/go-shorterurl/internal/app/logger"
	"github.com/kindenko/go-shorterurl/internal/app/zip"
)

var urls = make(map[string]string)

func main() {
	conf := config.NewCfg()
	newHandlers := handlers.NewHandlers(conf)

	r := chi.NewRouter()
	r.Use(logger.WithLogging)
	r.Use(zip.UnzipRequest)
	r.Use(zip.GzipHandler)

	r.Route("/", func(r chi.Router) {
		r.Post("/", newHandlers.PostHandler)
		r.Post("/api/shorten", newHandlers.PostJSONHandler)
		r.Get("/{shortUrl}", newHandlers.GetHandler)

	})

	log.Fatal(http.ListenAndServe(conf.Host, r))

}
