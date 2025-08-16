package store

import (
	"context"
)

type MemoryStorage struct {
	cache map[string]string
}

func (m *MemoryStorage) Ping(ctx context.Context) error {
	return nil
}

func (m *MemoryStorage) Get(ctx context.Context, code string) (value string, isFound bool) {
	value, ok := m.cache[code]
	return value, ok
}

func (m *MemoryStorage) Set(ctx context.Context, url ShortenedURL) (storedURL ShortenedURL, hasConflict bool, err error) {
	const key = "k"
	storedURLs, hasConflict, err := m.SetMany(ctx, map[string]ShortenedURL{key: url})
	if err != nil {
		return url, hasConflict, err
	}

	return storedURLs[key], hasConflict, err
}

func (m *MemoryStorage) SetMany(ctx context.Context, urls map[string]ShortenedURL) (storedURLs map[string]ShortenedURL, hasConflict bool, err error) {
	hasConflict = false
	storedURLs = make(map[string]ShortenedURL, len(urls))

	for key, url := range urls {
		for storedCode, storedURL := range m.cache {
			if storedURL == url.OriginalURL {
				hasConflict = true

				storedURLs[key] = ShortenedURL{
					Code:        storedCode,
					OriginalURL: storedURL,
				}
			}
		}
		m.cache[url.Code] = url.OriginalURL
	}
	return urls, hasConflict, nil
}

func (m *MemoryStorage) Close() error {
	return nil
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		cache: map[string]string{},
	}
}
