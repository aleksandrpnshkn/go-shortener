package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aleksandrpnshkn/go-shortener/internal/services/audit"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

// Ошибки Unshortener.
var (
	ErrShortURLWasDeleted = errors.New("short url was deleted")
)

// Unshortener - сервис для получения оригинальных URLов.
type Unshortener struct {
	urlsStorage       urls.Storage
	followedPublisher *audit.Publisher
}

// Unshorten получает оригинальный URL по short-коду из БД.
func (s *Unshortener) Unshorten(
	ctx context.Context,
	code types.Code,
	userID types.UserID,
) (types.OriginalURL, error) {
	followedAt := time.Now()

	url, err := s.urlsStorage.Get(ctx, code)
	if err != nil {
		return types.OriginalURL(""),
			fmt.Errorf("failed to get original url from storage: %w", err)
	}

	if url.IsDeleted {
		return types.OriginalURL(""), ErrShortURLWasDeleted
	}

	s.followedPublisher.Notify(ctx, audit.NewFollowEvent(
		followedAt,
		userID,
		url.OriginalURL,
	))

	return url.OriginalURL, nil
}

// NewUnshortener - создаёт новый Unshortener.
func NewUnshortener(
	urlsStorage urls.Storage,
	followedPublisher *audit.Publisher,
) *Unshortener {
	return &Unshortener{
		urlsStorage:       urlsStorage,
		followedPublisher: followedPublisher,
	}
}
