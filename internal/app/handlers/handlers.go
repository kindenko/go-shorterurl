package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kindenko/go-shorterurl/internal/app/structures"
)

type RequestJSON struct {
	URL string `json:"url"`
}

type ResponseJSON struct {
	Result string `json:"result"`
}

func (h *Handlers) PostHandler(w http.ResponseWriter, r *http.Request) {

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

	shortURL, err := h.storage.Save(url)
	if err != nil {
		fmt.Println(err)
	}

	resp := h.cfg.ResultURL + "/" + shortURL
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte(resp))
}

func (h *Handlers) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id := r.URL.Path[1:]
		url, err := h.storage.Get(id)
		if err != nil {
			http.Error(w, "Bad URL", http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handlers) PostJSONHandler(w http.ResponseWriter, r *http.Request) {
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

		shortURL, err := h.storage.Save(url)
		if err != nil {
			fmt.Println(err)
		}

		result := ResponseJSON{Result: h.cfg.ResultURL + "/" + shortURL}

		resp, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		w.Write(resp)
	}
}

func (h *Handlers) Batch(w http.ResponseWriter, r *http.Request) {
	var batches []structures.BatchEntity
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Batch: failed to read from body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &batches)
	fmt.Println("Batch request body", batches)
	if err != nil {
		fmt.Println("Batch: failed to unmarshal request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := h.storage.Batch(batches)
	if err != nil {
		fmt.Println("Batch: failed to save to database")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response, err := json.Marshal(result)
	fmt.Println("Batch response", string(response))
	if err != nil {
		fmt.Println("Batch: failed to marshal response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)

}

func (h *Handlers) PingDataBase(w http.ResponseWriter, _ *http.Request) {

	db, err := sql.Open("pgx", h.cfg.DataBaseString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
