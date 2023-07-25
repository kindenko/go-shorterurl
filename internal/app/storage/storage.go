package storage

import (
	"github.com/kindenko/go-shorterurl/config"
	"github.com/kindenko/go-shorterurl/internal/app/database"
	"github.com/kindenko/go-shorterurl/internal/app/structures"
)

type MyStorage interface {
	Save(fullURL string) (string, error)
	Get(shortURL string) (string, error)
	Batch(entities []structures.BatchEntity) ([]structures.BatchEntity, error)
}

type storage struct {
	defaultStorage MyStorage
}

func Init(cfg *config.AppConfig) MyStorage {
	var s storage

	switch {
	case cfg.DataBaseString != "":
		{
			s.defaultStorage = database.InitDB(cfg.DataBaseString, cfg.ResultURL)
			return &s
		}
	case cfg.FilePATH != "/tmp/short-url-db.json":
		{
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

func (s *storage) Batch(entities []structures.BatchEntity) ([]structures.BatchEntity, error) {
	return s.defaultStorage.Batch(entities)
}
