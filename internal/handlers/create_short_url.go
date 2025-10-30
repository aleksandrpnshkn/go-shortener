package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/aleksandrpnshkn/go-shortener/internal/services"
	"go.uber.org/zap"
)

func CreateShortURLPlain(
	shortener *services.Shortener,
	logger *zap.Logger,
	auther services.Auther,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain")

		url, err := io.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		defer req.Body.Close()

		user, err := auther.FromUserContext(req.Context())
		if err != nil {
			logger.Error("failed to get user", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		shortURL, hasConflict, err := shortener.Shorten(req.Context(), services.OriginalURL(url), user)
		if err != nil {
			logger.Error("failed to create plain short url", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		if hasConflict {
			res.WriteHeader(http.StatusConflict)
		} else {
			res.WriteHeader(http.StatusCreated)
		}

		res.Write([]byte(shortURL))
	}
}

type createShortURLRequest struct {
	URL string `json:"url"`
}

type createShortURLResponse struct {
	Result string `json:"result"`
}

func CreateShortURL(
	shortener *services.Shortener,
	logger *zap.Logger,
	auther services.Auther,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json")

		rawRequestData, err := io.ReadAll(req.Body)
		if err != nil {
			writeBadRequestError(res)
			return
		}
		defer req.Body.Close()

		var requestData createShortURLRequest
		err = json.Unmarshal(rawRequestData, &requestData)
		if err != nil {
			writeBadRequestError(res)
			return
		}

		if len(requestData.URL) == 0 {
			writeBadRequestError(res)
			return
		}

		user, err := auther.FromUserContext(req.Context())
		if err != nil {
			logger.Error("failed to get user", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		shortURL, hasConflict, err := shortener.Shorten(req.Context(), services.OriginalURL(requestData.URL), user)
		if err != nil {
			logger.Error("failed to create short url", zap.Error(err))
			writeInternalServerError(res)
			return
		}

		responseData := createShortURLResponse{
			Result: string(shortURL),
		}
		rawResponseData, err := json.Marshal(responseData)
		if err != nil {
			writeBadRequestError(res)
			return
		}

		if hasConflict {
			res.WriteHeader(http.StatusConflict)
		} else {
			res.WriteHeader(http.StatusCreated)
		}

		res.Write(rawResponseData)
	}
}
