package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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
	r.Use(zip.MiddlewareCompressGzip)
	r.Use(middleware.Compress(5, "text/html",
		"application/x-gzip",
		"text/plain",
		"application/json"))

	r.Route("/", func(r chi.Router) {
		r.Post("/", newHandlers.PostHandler)
		r.Post("/api/shorten", newHandlers.PostJSONHandler)
		r.Get("/{shortUrl}", newHandlers.GetHandler)
		r.Get("/ping", newHandlers.PingDataBase)
	})

	log.Fatal(http.ListenAndServe(conf.Host, r))

}
