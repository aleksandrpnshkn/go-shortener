package urls

import (
	"context"
	"errors"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

type ShortenedURL struct {
	Code        types.Code
	OriginalURL types.OriginalURL
	IsDeleted   bool
}

type DeleteCode struct {
	User users.User
	Code types.Code
}

type Storage interface {
	Ping(ctx context.Context) error

	Set(ctx context.Context, url ShortenedURL, user *users.User) (storedURL ShortenedURL, hasConflict bool, err error)

	SetMany(ctx context.Context, urls map[string]ShortenedURL, user *users.User) (storedURLs map[string]ShortenedURL, hasConflict bool, err error)

	Get(ctx context.Context, code types.Code) (ShortenedURL, error)

	GetByUserID(ctx context.Context, user *users.User) ([]ShortenedURL, error)

	DeleteManyByUserID(ctx context.Context, batch []DeleteCode) error

	Close() error
}

var (
	ErrCodeNotFound        = errors.New("code not found")
	ErrOriginalURLNotFound = errors.New("original url not found")
)
