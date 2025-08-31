package services

import (
	"context"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

type ShortenedURL struct {
	OriginalURL types.OriginalURL
	ShortURL    types.ShortURL
}

type Shortener struct {
	codesReserver CodesReserver
	urlsStorage   urls.Storage
	baseURL       string
}

func (s *Shortener) Shorten(ctx context.Context, originalURL types.OriginalURL, user *users.User) (shortURL types.ShortURL, hasConflict bool, err error) {
	code, err := s.codesReserver.GetCode(ctx)
	if err != nil {
		return shortURL, hasConflict, err
	}

	urlToStore := urls.ShortenedURL{
		Code:        code,
		OriginalURL: originalURL,
	}

	storedURL, hasConflict, err := s.urlsStorage.Set(ctx, urlToStore, user)
	if err != nil {
		return shortURL, hasConflict, err
	}

	shortURL = s.makeShortURL(storedURL.Code)

	return shortURL, hasConflict, nil
}

func (s *Shortener) ShortenMany(
	ctx context.Context,
	originalURLs map[string]types.OriginalURL,
	user *users.User,
) (shortURLs map[string]types.ShortURL, err error) {
	urlsToStore := make(map[string]urls.ShortenedURL, len(originalURLs))
	shortURLs = make(map[string]types.ShortURL, len(originalURLs))

	for correlationID, url := range originalURLs {
		code, err := s.codesReserver.GetCode(ctx)
		if err != nil {
			return shortURLs, err
		}

		urlsToStore[correlationID] = urls.ShortenedURL{
			Code:        code,
			OriginalURL: url,
		}
	}

	storedURLs, _, err := s.urlsStorage.SetMany(ctx, urlsToStore, user)
	if err != nil {
		return nil, err
	}

	for correlationID, url := range storedURLs {
		shortURLs[correlationID] = s.makeShortURL(url.Code)
	}

	return shortURLs, nil
}

func (s *Shortener) GetUserURLs(ctx context.Context, user *users.User) ([]ShortenedURL, error) {
	userURLs, err := s.urlsStorage.GetByUserID(ctx, user)
	if err != nil {
		return nil, err
	}

	result := []ShortenedURL{}

	for _, url := range userURLs {
		result = append(result, ShortenedURL{
			ShortURL:    s.makeShortURL(url.Code),
			OriginalURL: types.OriginalURL(url.OriginalURL),
		})
	}

	return result, nil
}

func (s *Shortener) makeShortURL(code types.Code) types.ShortURL {
	return types.ShortURL(s.baseURL + "/" + string(code))
}

func NewShortener(
	ctx context.Context,
	codesReserver CodesReserver,
	urlsStorage urls.Storage,
	baseURL string,
) *Shortener {
	shortener := Shortener{
		codesReserver: codesReserver,
		urlsStorage:   urlsStorage,
		baseURL:       baseURL,
	}

	return &shortener
}
