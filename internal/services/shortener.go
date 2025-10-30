package services

import (
	"context"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
)

type OriginalURL string

type ShortURL string

type ShortenedURL struct {
	OriginalURL OriginalURL
	ShortURL    ShortURL
}

type Shortener struct {
	codeGenerator CodeGenerator
	urlsStorage   urls.Storage
	baseURL       string
}

func (s *Shortener) Shorten(ctx context.Context, originalURL OriginalURL, user *users.User) (shortURL ShortURL, hasConflict bool, err error) {
	const fakeCorrelationID = "fake_id"

	shortURLs, hasConflict, err := s.ShortenMany(ctx, map[string]OriginalURL{fakeCorrelationID: originalURL}, user)
	if err != nil {
		return "", hasConflict, err
	}

	return shortURLs[fakeCorrelationID], hasConflict, nil
}

func (s *Shortener) ShortenMany(
	ctx context.Context,
	originalURLs map[string]OriginalURL,
	user *users.User,
) (shortURLs map[string]ShortURL, hasConflict bool, err error) {
	codesInBatch := make(map[Code]bool, len(originalURLs))
	urlsToStore := make(map[string]urls.ShortenedURL, len(originalURLs))
	shortURLs = make(map[string]ShortURL, len(originalURLs))

	for correlationID, url := range originalURLs {
		var code Code
		codeExistsInCurrentBatch := true
		codeExistsInDatabase := true
		for codeExistsInCurrentBatch || codeExistsInDatabase {
			code = s.codeGenerator.Generate()

			_, codeExistsInCurrentBatch = codesInBatch[code]

			_, err = s.urlsStorage.Get(ctx, string(code))
			codeExistsInDatabase = err == nil
		}

		codesInBatch[code] = true

		urlsToStore[correlationID] = urls.ShortenedURL{
			Code:        string(code),
			OriginalURL: string(url),
		}
	}

	storedURLs, hasConflict, err := s.urlsStorage.SetMany(ctx, urlsToStore, user)
	if err != nil {
		return nil, hasConflict, err
	}

	for correlationID, url := range storedURLs {
		shortURLs[correlationID] = s.makeShortURL(url.Code)
	}

	return shortURLs, hasConflict, nil
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
			OriginalURL: OriginalURL(url.OriginalURL),
		})
	}

	return result, nil
}

func (s *Shortener) makeShortURL(code string) ShortURL {
	return ShortURL(s.baseURL + "/" + code)
}

func NewShortener(
	codeGenerator CodeGenerator,
	urlsStorage urls.Storage,
	baseURL string,
) *Shortener {
	shortener := Shortener{
		codeGenerator,
		urlsStorage,
		baseURL,
	}

	return &shortener
}
