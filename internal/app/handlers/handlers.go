package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/kindenko/go-shorterurl/internal/app/storage"
)

var urls = make(map[string]string)

func (a *Handlers) PostHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusBadRequest)
		return
	}
	if string(body) == "" {
		http.Error(w, "Empty body!", http.StatusBadRequest)
		return
	}
	id := storage.RandString()
	urls[id] = string(body)
	resp := a.cfg.ResultURL + "/" + id
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte(resp))
}

func (a *Handlers) GetHandler(w http.ResponseWriter, r *http.Request) {
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
