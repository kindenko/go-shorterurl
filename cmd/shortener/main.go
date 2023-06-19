package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kindenko/go-shorterurl.git/config"
)

var urls = make(map[string]string)

func main() {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", postHandler)
		r.Get("/{shortUrl}", getHandler)
	})
	flag.Parse()
	log.Fatal(http.ListenAndServe(config.Config.Host, r))

}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusBadRequest)
			return
		}
		if string(body) == "" {
			http.Error(w, "Empty body!", http.StatusBadRequest)
			return
		}
		id := randString()
		urls[id] = string(body)
		resp := config.Config.ResultURL + "/" + id
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)

		w.Write([]byte(resp))
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id := r.URL.Path[1:]
		url, ok := urls[id]
		if !ok {
			http.Error(w, "Bad URL", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func randString() string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)

}
