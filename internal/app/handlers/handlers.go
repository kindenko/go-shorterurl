package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kindenko/go-shorterurl/internal/app/auth"
	e "github.com/kindenko/go-shorterurl/internal/app/errors"
	"github.com/kindenko/go-shorterurl/internal/app/structures"
	"github.com/kindenko/go-shorterurl/internal/app/utils"
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

	userID, err := auth.GetUserToken(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	url := string(body)
	short := utils.RandString(url)

	shortURL, err := h.storage.Save(url, short, userID)
	if err == e.ErrUniqueValue {
		resp := h.cfg.ResultURL + "/" + shortURL
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)

		w.Write([]byte(resp))

		return
	}

	resp := h.cfg.ResultURL + "/" + shortURL
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte(resp))
}

func (h *Handlers) GetHandler(w http.ResponseWriter, r *http.Request) {

	log.Println(h.cfg.ResultURL)
	if r.Method == http.MethodGet {
		id := r.URL.Path[1:]
		url, isDeleted, err := h.storage.Get(id)
		if err != nil {
			http.Error(w, "Bad URL", http.StatusBadRequest)
			return
		}

		if isDeleted == 1 {
			w.WriteHeader(http.StatusGone)
			return
		}

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handlers) GetUsersURLs(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserToken(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	batch, err := h.storage.GetBatchByUserID(userID)
	if err != nil {
		log.Println("Failed to fetch user data")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response, err := json.Marshal(batch)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
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
		short := utils.RandString(url)

		userID, err := auth.GetUserToken(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		shortURL, err := h.storage.Save(url, short, userID)
		if err == e.ErrUniqueValue {
			result := ResponseJSON{Result: h.cfg.ResultURL + "/" + shortURL}
			resp, err := json.Marshal(result)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)

			w.Write(resp)

			return
		}
		if err != nil {
			log.Println(err)
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
	userID, err := auth.GetUserToken(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var batches []structures.BatchEntity
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Batch: failed to read from body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &batches)
	log.Println("Batch request body", batches)
	if err != nil {
		log.Println("Batch: failed to unmarshal request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := h.storage.Batch(batches, userID)
	if err != nil {
		log.Println("Batch: failed to save to database")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response, err := json.Marshal(result)
	log.Println("Batch response", string(response))
	if err != nil {
		log.Println("Batch: failed to marshal response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)

}

func (h *Handlers) DeleteURLs(w http.ResponseWriter, r *http.Request) {
	var inputArray []string
	userID, err := auth.GetUserToken(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	urls, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("DeleteURLs: failed to read body, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(urls, &inputArray)
	if err != nil {
		log.Printf("DeleteURLs: failed to unmarshal input request, %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputCh := addShortURLs(inputArray)
	go h.MarkAsDeleted(inputCh, userID)

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handlers) MarkAsDeleted(inputShort chan string, userID string) {
	for v := range inputShort {
		err := h.storage.DeleteByUserIDAndShort(userID, v)
		if err != nil {
			log.Print("Failed to mark deleted by short")
		}
	}
}

func addShortURLs(input []string) chan string {
	inputCh := make(chan string, 10)

	go func() {
		defer close(inputCh)
		for _, url := range input {
			inputCh <- url
		}
	}()

	return inputCh
}

func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	if err := h.storage.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
