package web

import (
	"log"
	"net/http"

	"github.com/kindenko/go-shorterurl/config"
	"github.com/kindenko/go-shorterurl/internal/app/storage"
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
	s.srv.GET("/:shortUrl", s.GetOriginalURL)

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

	return nil
}
