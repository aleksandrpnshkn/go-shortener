package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aleksandrpnshkn/go-shortener/internal/services/audit"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

var (
	ErrShortURLWasDeleted = errors.New("short url was deleted")
)

type Unshortener struct {
	urlsStorage       urls.Storage
	followedPublisher *audit.Publisher
}

func (s *Unshortener) Unshorten(
	ctx context.Context,
	code types.Code,
	user *users.User,
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
		user,
		url.OriginalURL,
	))

	return url.OriginalURL, nil
}

func NewUnshortener(
	urlsStorage urls.Storage,
	followedPublisher *audit.Publisher,
) *Unshortener {
	return &Unshortener{
		urlsStorage:       urlsStorage,
		followedPublisher: followedPublisher,
	}
}
