package web

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/kindenko/go-shorterurl/config"
	"github.com/kindenko/go-shorterurl/internal/app/storage"
	"github.com/kindenko/go-shorterurl/internal/app/utils"
	"github.com/labstack/echo/v4"
)

type Server struct {
	store storage.MyStorage
	cfg   config.AppConfig
	srv   *echo.Echo
}

func NewServer(cfg config.AppConfig) *Server {
	return &Server{
		store: storage.Init(&cfg),
		cfg:   cfg,
		srv:   echo.New(),
	}
}

func (s *Server) Run() error {
	s.srv.Get("/:shortUrl", s.GetOriginalURL)

	return s.srv.Start(s.cfg.Host)
}

func (h *Server) GetOriginalURL(c echo.Context) error {

	log.Println(h.cfg.ResultURL)
	id := c.Request().URL.Path[1:]
	url, err := h.store.Get(id)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad URL")
	}

	c.Response().Header().Set("Location", url)
	c.Response().WriteHeader(http.StatusTemporaryRedirect)
}

func (s *Server) PostJSONHandler(w http.ResponseWriter, r *http.Request) {
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

		shortURL, err := h.storage.Save(url, short)
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
