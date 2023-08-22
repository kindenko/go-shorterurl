package storage

import (
	"github.com/kindenko/go-shorterurl/internal/app/structures"
)

type MemoryStorage struct {
	store map[string]string
}

func InitMemory() *MemoryStorage {
	return &MemoryStorage{
		store: make(map[string]string),
	}
}

func (m *MemoryStorage) Save(fullURL string, shortURL string) (string, error) {
	m.store[shortURL] = fullURL

	return shortURL, nil
}

func (m *MemoryStorage) Get(shortURL string) (string, error) {
	long, ok := m.store[shortURL]

	if ok {
		return long, nil
	}
	return "Missing url", nil
}

func (m *MemoryStorage) Batch(entities []structures.BatchEntity) ([]structures.BatchEntity, error) {
	panic("Missing method")
}
