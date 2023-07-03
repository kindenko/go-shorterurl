package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kindenko/go-shorterurl/internal/app/storage"
)

type RequestJSON struct {
	URL string `json:"url"`
}

type ResponseJSON struct {
	Result string `json:"result"`
}

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
	url := string(body)
	id := storage.RandString()
	urls[id] = string(body)
	resp := a.cfg.ResultURL + "/" + id

	fileStorage := storage.NewFileStorage()

	fileStorage.Short = id
	fileStorage.Original = url

	fmt.Println(a.cfg.FilePATH)
	err = storage.SaveToFile(fileStorage, a.cfg.FilePATH)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

func (a *Handlers) PostJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req RequestJSON
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		url := string(req.URL)
		id := storage.RandString()
		urls[id] = string(req.URL)

		result := ResponseJSON{Result: a.cfg.ResultURL + "/" + id}

		resp, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		fileStorage := storage.NewFileStorage()

		fileStorage.Short = id
		fileStorage.Original = url

		err = storage.SaveToFile(fileStorage, a.cfg.FilePATH)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		w.Write(resp)
	}
}
