package urls

import (
	"context"
	"errors"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
)

type ShortenedURL struct {
	Code        string
	OriginalURL string
}

type Storage interface {
	Ping(ctx context.Context) error

	Set(ctx context.Context, url ShortenedURL, user *users.User) (storedURL ShortenedURL, hasConflict bool, err error)

	SetMany(ctx context.Context, urls map[string]ShortenedURL, user *users.User) (storedURLs map[string]ShortenedURL, hasConflict bool, err error)

	Get(ctx context.Context, code string) (originalURL string, err error)

	GetByUserID(ctx context.Context, user *users.User) ([]ShortenedURL, error)

	Close() error
}

var (
	ErrCodeNotFound        = errors.New("code not found")
	ErrOriginalURLNotFound = errors.New("original url not found")
)
