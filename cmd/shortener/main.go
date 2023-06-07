package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

var urls = make(map[string]string)

func mainPost(w http.ResponseWriter, r *http.Request) {
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
	} else if r.Method == http.MethodGet {
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

func main() {
	urls = make(map[string]string)

	mux := http.NewServeMux()

	mux.HandleFunc("/", mainPost)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
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
