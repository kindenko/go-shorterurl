package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var urls = make(map[string]string)

func main() {
	urls = make(map[string]string)

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", newPost)
		r.Get("/{shortUrl}", newGet)
	})

	log.Fatal(http.ListenAndServe(":8080", r))

}

func newPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Err: %s", err), http.StatusBadRequest)
			return
		}
		if string(b) == "" {
			http.Error(w, "Empty body!", http.StatusBadRequest)
			return
		}
		url := string(b)
		if url == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id := randstring()
		urls[id] = string(b)
		resp := "http://localhost:8080/" + id
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)

		w.Write([]byte(resp))
	}
}

func newGet(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id := r.URL.Path[1:]
		url, ok := urls[id]
		if !ok {
			http.Error(w, "Bad URL", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func randstring() string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)

}
