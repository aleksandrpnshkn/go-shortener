package store

import "context"

type ShortenedURL struct {
	Code        string
	OriginalURL string
}

type Storage interface {
	Set(ctx context.Context, url ShortenedURL) error

	SetMany(ctx context.Context, urls []ShortenedURL) error

	Get(ctx context.Context, code string) (originalURL string, isFound bool)
}
