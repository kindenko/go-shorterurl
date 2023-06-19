package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kindenko/go-shorterurl.git/config"
	"github.com/kindenko/go-shorterurl.git/internal/app/handlers"
)

var urls = make(map[string]string)

func main() {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.PostHandler)
		r.Get("/{shortUrl}", handlers.GetHandler)
	})
	flag.Parse()
	log.Fatal(http.ListenAndServe(config.Config.Host, r))

}
