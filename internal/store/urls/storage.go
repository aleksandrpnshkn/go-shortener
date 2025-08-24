package urls

import (
	"context"
	"errors"
)

type ShortenedURL struct {
	Code        string
	OriginalURL string
}

type Storage interface {
	Ping(ctx context.Context) error

	Set(ctx context.Context, url ShortenedURL) (storedURL ShortenedURL, hasConflict bool, err error)

	SetMany(ctx context.Context, urls map[string]ShortenedURL) (storedURLs map[string]ShortenedURL, hasConflict bool, err error)

	Get(ctx context.Context, code string) (originalURL string, err error)

	Close() error
}

var (
	ErrCodeNotFound        = errors.New("code not found")
	ErrOriginalURLNotFound = errors.New("original url not found")
)
