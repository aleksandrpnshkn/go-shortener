package urls

import (
	"context"

	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

type MemoryStorage struct {
	cache map[types.Code]ShortenedURL
}

func (m *MemoryStorage) Ping(ctx context.Context) error {
	return nil
}

func (m *MemoryStorage) Get(ctx context.Context, code types.Code) (ShortenedURL, error) {
	value, ok := m.cache[code]
	if !ok {
		return ShortenedURL{}, ErrCodeNotFound
	}
	return value, nil
}

func (m *MemoryStorage) GetByUserID(ctx context.Context, userID types.UserID) ([]ShortenedURL, error) {
	return []ShortenedURL{}, nil
}

func (m *MemoryStorage) DeleteManyByUserID(ctx context.Context, commands []DeleteCode) error {
	return nil
}

func (m *MemoryStorage) Set(ctx context.Context, url ShortenedURL, userID types.UserID) (storedURL ShortenedURL, hasConflict bool, err error) {
	const key = "k"
	storedURLs, hasConflict, err := m.SetMany(ctx, map[string]ShortenedURL{key: url}, userID)
	if err != nil {
		return url, hasConflict, err
	}

	return storedURLs[key], hasConflict, err
}

func (m *MemoryStorage) SetMany(ctx context.Context, urls map[string]ShortenedURL, userID types.UserID) (storedURLs map[string]ShortenedURL, hasConflicts bool, err error) {
	hasConflicts = false
	storedURLs = make(map[string]ShortenedURL, len(urls))

	for key, url := range urls {
		for _, storedURL := range m.cache {
			if storedURL.OriginalURL == url.OriginalURL {
				hasConflicts = true

				storedURLs[key] = storedURL
			}
		}
		m.cache[url.Code] = url
	}
	return urls, hasConflicts, nil
}

func (m *MemoryStorage) Close() error {
	return nil
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		cache: map[types.Code]ShortenedURL{},
	}
}
