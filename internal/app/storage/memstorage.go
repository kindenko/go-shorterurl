package storage

import (
	"github.com/kindenko/go-shorterurl/internal/app/utils"
)

type MemoryStorage struct {
	store map[string]string
}

func InitMemory() *MemoryStorage {
	return &MemoryStorage{
		store: make(map[string]string),
	}
}

func (m *MemoryStorage) Save(fullUrl string) (string, error) {
	shortUrl := utils.RandString()
	m.store[shortUrl] = fullUrl

	return shortUrl, nil
}

func (m *MemoryStorage) Get(shortURL string) (string, error) {
	long, ok := m.store[shortURL]

	if ok {
		return long, nil
	}
	return "Missing url", nil
}
