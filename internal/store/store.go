package store

type Storage interface {
	Set(shortURL string, originalURL string) error

	Get(shortURL string) (originalURL string, isFound bool)
}
