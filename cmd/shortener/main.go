package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kindenko/go-shorterurl/config"
	"github.com/kindenko/go-shorterurl/internal/app/handlers"
)

var urls = make(map[string]string)

func main() {
	conf := config.NewCfg()
	newHandlers := handlers.NewHandlers(conf, chi.NewRouter())

	log.Fatal(http.ListenAndServe(conf.Host, newHandlers))
}
