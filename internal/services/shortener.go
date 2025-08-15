package services

import (
	"context"
	"maps"
	"slices"

	"github.com/aleksandrpnshkn/go-shortener/internal/store"
)

type OriginalURL string

type ShortURL string

type Shortener struct {
	codeGenerator CodeGenerator
	urlsStorage   store.Storage
	baseURL       string
}

func (s *Shortener) Shorten(ctx context.Context, URL OriginalURL) (ShortURL, error) {
	const fakeCorrelationID = "fake_id"

	shortURLs, err := s.ShortenMany(ctx, map[string]OriginalURL{fakeCorrelationID: URL})
	if err != nil {
		return "", err
	}

	return shortURLs[fakeCorrelationID], nil
}

func (s *Shortener) ShortenMany(ctx context.Context, URLs map[string]OriginalURL) (map[string]ShortURL, error) {
	URLsToStore := make(map[Code]store.ShortenedURL, len(URLs))
	shortURLs := make(map[string]ShortURL, len(URLs))

	for correlationID, URL := range URLs {
		var code Code
		codeExistsInCurrentBatch := true
		codeExistsInDatabase := true
		for codeExistsInCurrentBatch || codeExistsInDatabase {
			code = s.codeGenerator.Generate()

			_, codeExistsInCurrentBatch = URLsToStore[code]
			_, codeExistsInDatabase = s.urlsStorage.Get(ctx, string(code))
		}

		URLsToStore[code] = store.ShortenedURL{
			Code:        string(code),
			OriginalURL: string(URL),
		}

		shortURLs[correlationID] = ShortURL(s.baseURL + "/" + string(code))
	}

	err := s.urlsStorage.SetMany(ctx, slices.Collect(maps.Values(URLsToStore)))
	if err != nil {
		return nil, err
	}

	return shortURLs, nil
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
