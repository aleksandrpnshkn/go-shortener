package store

import "context"

type MemoryStorage struct {
	cache map[string]string
}

func (m *MemoryStorage) Get(ctx context.Context, code string) (value string, isFound bool) {
	value, ok := m.cache[code]
	return value, ok
}

func (m *MemoryStorage) Set(ctx context.Context, url ShortenedURL) error {
	m.cache[url.Code] = url.OriginalURL
	return nil
}

func (m *MemoryStorage) SetMany(ctx context.Context, urls []ShortenedURL) error {
	for _, url := range urls {
		m.Set(ctx, url)
	}
	return nil
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		cache: map[string]string{},
	}
}
