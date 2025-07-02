package services

import "github.com/aleksandrpnshkn/go-shortener/internal/store"

type FullURL string

type FullURLsStorage struct {
	store store.Store
}

func (s *FullURLsStorage) Set(code Code, url FullURL) error {
	return s.store.Set(string(code), string(url))
}

func (s *FullURLsStorage) Get(code Code) (url FullURL, isFound bool) {
	value, isFound := s.store.Get(string(code))
	return FullURL(value), isFound
}

func NewFullURLsStorage(store store.Store) *FullURLsStorage {
	fullURLsStorage := FullURLsStorage{
		store: store,
	}

	return &fullURLsStorage
}

func NewFullURLsTestStorage() *FullURLsStorage {
	return NewFullURLsStorage(store.NewMemoryStore())
}
