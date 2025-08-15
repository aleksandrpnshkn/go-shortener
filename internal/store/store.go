package store

import "context"

type Storage interface {
	Set(ctx context.Context, shortURL string, originalURL string) error

	Get(ctx context.Context, shortURL string) (originalURL string, isFound bool)
}
