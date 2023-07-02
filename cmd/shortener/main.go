package main

import (
	"flag"
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
	r.Route("/", func(r chi.Router) {
		r.Post("/", logger.WithLogging(zip.GzipMiddleware(newHandlers.PostHandler)))
		r.Post("/api/shorten", logger.WithLogging(zip.GzipMiddleware(newHandlers.PostJSONHandler)))
		r.Get("/{shortUrl}", logger.WithLogging(zip.GzipMiddleware(newHandlers.GetHandler)))
	})
	flag.Parse()
	log.Fatal(http.ListenAndServe(conf.Host, r))

}
