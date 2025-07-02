package services

import "github.com/aleksandrpnshkn/go-shortener/internal/store"

type OriginalURL string

type URLsStorage struct {
	storage store.Storage
}

func (s *URLsStorage) Set(code Code, url OriginalURL) error {
	return s.storage.Set(string(code), string(url))
}

func (s *URLsStorage) Get(code Code) (url OriginalURL, isFound bool) {
	originalURL, isFound := s.storage.Get(string(code))
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
