package services

import (
	"context"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
)

type OriginalURL string

type ShortURL string

type Shortener struct {
	codeGenerator CodeGenerator
	urlsStorage   urls.Storage
	baseURL       string
}

func (s *Shortener) Shorten(ctx context.Context, originalURL OriginalURL) (shortURL ShortURL, hasConflict bool, err error) {
	const fakeCorrelationID = "fake_id"

	shortURLs, hasConflict, err := s.ShortenMany(ctx, map[string]OriginalURL{fakeCorrelationID: originalURL})
	if err != nil {
		return "", hasConflict, err
	}

	return shortURLs[fakeCorrelationID], hasConflict, nil
}

func (s *Shortener) ShortenMany(ctx context.Context, originalURLs map[string]OriginalURL) (shortURLs map[string]ShortURL, hasConflict bool, err error) {
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

	storedURLs, hasConflict, err := s.urlsStorage.SetMany(ctx, urlsToStore)
	if err != nil {
		return nil, hasConflict, err
	}

	for correlationID, url := range storedURLs {
		shortURLs[correlationID] = ShortURL(s.baseURL + "/" + url.Code)
	}

	return shortURLs, hasConflict, nil
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
