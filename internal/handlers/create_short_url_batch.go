package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/aleksandrpnshkn/go-shortener/internal/services"
	"go.uber.org/zap"
)

type originalURL struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type shortURL struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func CreateShortURLBatch(
	shortener *services.Shortener,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json")

		rawRequestData, err := io.ReadAll(req.Body)
		if err != nil {
			writeBadRequestError(res)
			return
		}
		defer req.Body.Close()

		var requestData []originalURL
		err = json.Unmarshal(rawRequestData, &requestData)
		if err != nil {
			writeBadRequestError(res)
			return
		}

		if len(requestData) == 0 {
			writeBadRequestError(res)
			return
		}

		URLs := make(map[string]services.OriginalURL, len(requestData))

		for _, url := range requestData {
			URLs[url.CorrelationID] = services.OriginalURL(url.OriginalURL)
		}

		shortURLs, _, err := shortener.ShortenMany(req.Context(), URLs)
		if err != nil {
			logger.Error("failed to create short url", zap.Error(err))
			writeInternalServerError(res)
			return
		}

		result := []shortURL{}

		for correlationID, url := range shortURLs {
			result = append(result, shortURL{
				CorrelationID: correlationID,
				ShortURL:      string(url),
			})
		}

		rawResponseData, err := json.Marshal(result)
		if err != nil {
			writeBadRequestError(res)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write(rawResponseData)
	}
}
