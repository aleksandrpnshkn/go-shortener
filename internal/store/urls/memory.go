package urls

import (
	"context"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
)

type MemoryStorage struct {
	cache map[string]string
}

func (m *MemoryStorage) Ping(ctx context.Context) error {
	return nil
}

func (m *MemoryStorage) Get(ctx context.Context, code string) (value string, err error) {
	value, ok := m.cache[code]
	if !ok {
		return "", ErrCodeNotFound
	}
	return value, nil
}

func (m *MemoryStorage) GetByUserID(ctx context.Context, user *users.User) ([]ShortenedURL, error) {
	return []ShortenedURL{}, nil
}

func (m *MemoryStorage) Set(ctx context.Context, url ShortenedURL, user *users.User) (storedURL ShortenedURL, hasConflict bool, err error) {
	const key = "k"
	storedURLs, hasConflict, err := m.SetMany(ctx, map[string]ShortenedURL{key: url}, user)
	if err != nil {
		return url, hasConflict, err
	}

	return storedURLs[key], hasConflict, err
}

func (m *MemoryStorage) SetMany(ctx context.Context, urls map[string]ShortenedURL, user *users.User) (storedURLs map[string]ShortenedURL, hasConflict bool, err error) {
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
