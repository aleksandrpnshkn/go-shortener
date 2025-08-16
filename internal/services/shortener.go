package services

import (
	"context"

	"github.com/aleksandrpnshkn/go-shortener/internal/store"
)

type OriginalURL string

type ShortURL string

type Shortener struct {
	codeGenerator CodeGenerator
	urlsStorage   store.Storage
	baseURL       string
}

func (s *Shortener) Shorten(ctx context.Context, URL OriginalURL) (shortURL ShortURL, hasConflict bool, err error) {
	const fakeCorrelationID = "fake_id"

	shortURLs, hasConflict, err := s.ShortenMany(ctx, map[string]OriginalURL{fakeCorrelationID: URL})
	if err != nil {
		return "", hasConflict, err
	}

	return shortURLs[fakeCorrelationID], hasConflict, nil
}

func (s *Shortener) ShortenMany(ctx context.Context, URLs map[string]OriginalURL) (shortURLs map[string]ShortURL, hasConflict bool, err error) {
	codesInBatch := make(map[Code]bool, len(URLs))
	URLsToStore := make(map[string]store.ShortenedURL, len(URLs))
	shortURLs = make(map[string]ShortURL, len(URLs))

	for correlationID, URL := range URLs {
		var code Code
		codeExistsInCurrentBatch := true
		codeExistsInDatabase := true
		for codeExistsInCurrentBatch || codeExistsInDatabase {
			code = s.codeGenerator.Generate()

			_, codeExistsInCurrentBatch = codesInBatch[code]
			_, codeExistsInDatabase = s.urlsStorage.Get(ctx, string(code))
		}

		codesInBatch[code] = true

		URLsToStore[correlationID] = store.ShortenedURL{
			Code:        string(code),
			OriginalURL: string(URL),
		}
	}

	storedURLs, hasConflict, err := s.urlsStorage.SetMany(ctx, URLsToStore)
	if err != nil {
		return nil, hasConflict, err
	}

	for correlationID, URL := range storedURLs {
		shortURLs[correlationID] = ShortURL(s.baseURL + "/" + URL.Code)
	}

	return shortURLs, hasConflict, nil
}

func NewShortener(
	codeGenerator CodeGenerator,
	urlsStorage store.Storage,
	baseURL string,
) *Shortener {
	shortener := Shortener{
		codeGenerator,
		urlsStorage,
		baseURL,
	}

	return &shortener
}
