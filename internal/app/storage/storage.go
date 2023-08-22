package storage

import (
	"log"

	"github.com/kindenko/go-shorterurl/config"
	"github.com/kindenko/go-shorterurl/internal/app/database"
	e "github.com/kindenko/go-shorterurl/internal/app/errors"
	"github.com/kindenko/go-shorterurl/internal/app/structures"
)

type MyStorage interface {
	Save(fullURL string, shortURL string) (string, error)
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
			log.Println("DB")
			s.defaultStorage = database.InitDB(cfg.DataBaseString, cfg.ResultURL)
			return &s
		}
	case cfg.FilePATH != "/tmp/short-url-db.json":
		{
			log.Println("FILE")
			s.defaultStorage = InitFileDB(cfg.FilePATH)
			return &s
		}
	}
	log.Println("MEM")
	s.defaultStorage = InitMemory()

	return &s
}

func (s *storage) Save(full string, short string) (string, error) {
	short, err := s.defaultStorage.Save(full, short)
	if err == e.ErrUniqueValue {
		return short, err
	}

	if err != nil && err == e.ErrUniqueValue {
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
