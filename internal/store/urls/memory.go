package urls

import (
	"context"

	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

// MemoryStorage - in-memory хранилище сокращённых ссылок.
// Реализует ограниченный набор методов.
type MemoryStorage struct {
	cache map[types.Code]ShortenedURL
}

// Ping проверяет доступность хранилища.
func (m *MemoryStorage) Ping(ctx context.Context) error {
	return nil
}

// Get достаёт сокращённую ссылку по её коду.
func (m *MemoryStorage) Get(ctx context.Context, code types.Code) (ShortenedURL, error) {
	value, ok := m.cache[code]
	if !ok {
		return ShortenedURL{}, ErrCodeNotFound
	}
	return value, nil
}

// GetByUserID - заглушка для доставания всех сокращённых ссылок пользователя.
func (m *MemoryStorage) GetByUserID(ctx context.Context, userID types.UserID) ([]ShortenedURL, error) {
	return []ShortenedURL{}, nil
}

// DeleteManyByUserID - заглушка для удаления всех сокращённых ссылок пользователя.
func (m *MemoryStorage) DeleteManyByUserID(ctx context.Context, commands []DeleteCode) error {
	return nil
}

// Set сохраняет короткую ссылку в хранилище.
// А также проверяет наличие дублей.
func (m *MemoryStorage) Set(ctx context.Context, url ShortenedURL, userID types.UserID) (storedURL ShortenedURL, hasConflict bool, err error) {
	const key = "k"
	storedURLs, hasConflict, err := m.SetMany(ctx, map[string]ShortenedURL{key: url}, userID)
	if err != nil {
		return url, hasConflict, err
	}

	return storedURLs[key], hasConflict, err
}

// SetMany сохраняет множество коротких ссылок в хранилищк.
// А также проверяет наличие дублей.
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

// Close заглушка для закрытия соединения с хранилищем.
func (m *MemoryStorage) Close() error {
	return nil
}

// NewFileStorage создаёт новое in memory хранилище для URLов.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		cache: map[types.Code]ShortenedURL{},
	}
}
