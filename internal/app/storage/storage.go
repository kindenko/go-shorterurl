package storage

import (
	"fmt"

	"github.com/kindenko/go-shorterurl/config"
	"github.com/kindenko/go-shorterurl/internal/app/database"
)

type MyStorage interface {
	Save(fullURL string) (string, error)
	Get(shortURL string) (string, error)
}

type storage struct {
	defaultStorage MyStorage
}

func Init(cfg *config.AppConfig) MyStorage {
	var s storage

	switch {
	case cfg.DataBaseString != "":
		{
			fmt.Println("db")
			s.defaultStorage = database.InitDB(cfg.DataBaseString, cfg.ResultURL)
			return &s
		}
	case cfg.FilePATH != "/tmp/short-url-db.json":
		{
			fmt.Println("FILE")
			s.defaultStorage = InitFileDB(cfg.FilePATH)
			return &s
		}
	}
	s.defaultStorage = InitMemory()

	return &s
}

func (s *storage) Save(full string) (string, error) {
	short, err := s.defaultStorage.Save(full)
	if err != nil {
		return "", err
	}
	return short, nil
}

func (s *storage) Get(short string) (string, error) {
	full, err := s.defaultStorage.Get(short)
	if err != nil {
		return "", err
	}
	return full, nil
}
