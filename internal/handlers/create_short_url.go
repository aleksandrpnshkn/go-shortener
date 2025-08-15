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
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain")

		URL, err := io.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		defer req.Body.Close()

		shortURL, err := shortener.Shorten(req.Context(), services.OriginalURL(URL))
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusCreated)
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

		shortURL, err := shortener.Shorten(req.Context(), services.OriginalURL(requestData.URL))
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

		res.WriteHeader(http.StatusCreated)
		res.Write(rawResponseData)
	}
}
