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
	codeGenerator CodeGenerator
	urlsStorage   urls.Storage
	baseURL       string
}

func (s *Shortener) Shorten(ctx context.Context, originalURL types.OriginalURL, user *users.User) (shortURL types.ShortURL, hasConflict bool, err error) {
	const fakeCorrelationID = "fake_id"

	shortURLs, hasConflict, err := s.ShortenMany(ctx, map[string]types.OriginalURL{fakeCorrelationID: originalURL}, user)
	if err != nil {
		return "", hasConflict, err
	}

	return shortURLs[fakeCorrelationID], hasConflict, nil
}

func (s *Shortener) ShortenMany(
	ctx context.Context,
	originalURLs map[string]types.OriginalURL,
	user *users.User,
) (shortURLs map[string]types.ShortURL, hasConflict bool, err error) {
	codesInBatch := make(map[types.Code]bool, len(originalURLs))
	urlsToStore := make(map[string]urls.ShortenedURL, len(originalURLs))
	shortURLs = make(map[string]types.ShortURL, len(originalURLs))

	for correlationID, url := range originalURLs {
		var code types.Code
		codeExistsInCurrentBatch := true
		codeExistsInDatabase := true
		for codeExistsInCurrentBatch || codeExistsInDatabase {
			code = s.codeGenerator.Generate()

			_, codeExistsInCurrentBatch = codesInBatch[code]

			_, err = s.urlsStorage.Get(ctx, code)
			codeExistsInDatabase = err == nil
		}

		codesInBatch[code] = true

		urlsToStore[correlationID] = urls.ShortenedURL{
			Code:        code,
			OriginalURL: url,
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
			OriginalURL: types.OriginalURL(url.OriginalURL),
		})
	}

	return result, nil
}

func (s *Shortener) DeleteUserURLs(ctx context.Context, codes []types.Code, user *users.User) {
	s.urlsStorage.DeleteManyByUserID(ctx, codes, user)
}

func (s *Shortener) makeShortURL(code types.Code) types.ShortURL {
	return types.ShortURL(s.baseURL + "/" + string(code))
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
