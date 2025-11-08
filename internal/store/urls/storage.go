package urls

import (
	"context"
	"errors"

	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

type ShortenedURL struct {
	Code        types.Code
	OriginalURL types.OriginalURL
	IsDeleted   bool
}

type DeleteCode struct {
	UserID types.UserID
	Code   types.Code
}

type Storage interface {
	Ping(ctx context.Context) error

	Set(ctx context.Context, url ShortenedURL, userID types.UserID) (storedURL ShortenedURL, hasConflict bool, err error)

	SetMany(ctx context.Context, urls map[string]ShortenedURL, userID types.UserID) (storedURLs map[string]ShortenedURL, hasConflicts bool, err error)

	Get(ctx context.Context, code types.Code) (ShortenedURL, error)

	GetByUserID(ctx context.Context, userID types.UserID) ([]ShortenedURL, error)

	DeleteManyByUserID(ctx context.Context, batch []DeleteCode) error

	Close() error
}

var (
	ErrCodeNotFound        = errors.New("code not found")
	ErrOriginalURLNotFound = errors.New("original url not found")
)
