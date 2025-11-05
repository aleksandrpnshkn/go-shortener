package urls

import (
	"context"
	"errors"

	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

// ShortenedURL - укороченная ссылка. Содержит в себе оригинал и доп. инфу.
type ShortenedURL struct {
	Code        types.Code
	OriginalURL types.OriginalURL
	IsDeleted   bool
}

// DeleteCode - вспомогательная структура для массового удаления коротких ссылок у разных пользователей.
type DeleteCode struct {
	UserID types.UserID
	Code   types.Code
}

// Storage - интерфейс для работы с хранилищем сокращённых ссылок.
type Storage interface {
	// Ping проверяет доступность БД.
	Ping(ctx context.Context) error

	// Set сохраняет короткую ссылку в БД.
	// А также проверяет наличие дублей.
	Set(ctx context.Context, url ShortenedURL, userID types.UserID) (storedURL ShortenedURL, hasConflict bool, err error)

	// SetMany сохраняет множество коротких ссылок в БД.
	// А также проверяет наличие дублей.
	SetMany(ctx context.Context, urls map[string]ShortenedURL, userID types.UserID) (storedURLs map[string]ShortenedURL, hasConflicts bool, err error)

	// Get достаёт сокращённую ссылку по её коду.
	Get(ctx context.Context, code types.Code) (ShortenedURL, error)

	// GetByUserID достаёт все сокращённые ссылки пользователя.
	GetByUserID(ctx context.Context, userID types.UserID) ([]ShortenedURL, error)

	// DeleteManyByUserID удаляет все сокращённые ссылки пользователя.
	DeleteManyByUserID(ctx context.Context, batch []DeleteCode) error

	// Close закрывает соединение с хранилищем, если этого требует реализация.
	// Вызывается при завершении работы программы.
	Close() error
}

var (
	ErrCodeNotFound        = errors.New("code not found")
	ErrOriginalURLNotFound = errors.New("original url not found")
)
