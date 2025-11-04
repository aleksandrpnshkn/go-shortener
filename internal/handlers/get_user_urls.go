package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aleksandrpnshkn/go-shortener/internal/services"
	"go.uber.org/zap"
)

type shortenedURL struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

func GetUserURLs(
	shortener *services.Shortener,
	logger *zap.Logger,
	auther services.Auther,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json")

		userID, err := auther.FromUserContext(req.Context())
		if err != nil {
			logger.Error("failed to get user", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		userURLs, err := shortener.GetUserURLs(req.Context(), userID)
		if err != nil {
			logger.Error("failed to get user urls",
				zap.Error(err),
				zap.Int64("user_id", int64(userID)),
			)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(userURLs) == 0 {
			res.WriteHeader(http.StatusNoContent)
			return
		}

		result := []shortenedURL{}

		for _, url := range userURLs {
			result = append(result, shortenedURL{
				ShortURL:    string(url.ShortURL),
				OriginalURL: string(url.OriginalURL),
			})
		}

		rawResponseData, err := json.Marshal(result)
		if err != nil {
			writeBadRequestError(res)
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Write(rawResponseData)
	}
}
