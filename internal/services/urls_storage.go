package services

import (
	"context"

	"github.com/aleksandrpnshkn/go-shortener/internal/store"
)

type OriginalURL string

type URLsStorage struct {
	storage store.Storage
}

func (s *URLsStorage) Set(ctx context.Context, code Code, url OriginalURL) error {
	return s.storage.Set(ctx, string(code), string(url))
}

func (s *URLsStorage) Get(ctx context.Context, code Code) (url OriginalURL, isFound bool) {
	originalURL, isFound := s.storage.Get(ctx, string(code))
	return OriginalURL(originalURL), isFound
}

func NewURLsStorage(storage store.Storage) *URLsStorage {
	return &URLsStorage{
		storage: storage,
	}
}

func NewURLsTestStorage() *URLsStorage {
	return NewURLsStorage(store.NewMemoryStorage())
}
