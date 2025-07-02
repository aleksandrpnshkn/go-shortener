package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aleksandrpnshkn/go-shortener/internal/services"
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

		shortURL := shortener.Shorten(services.OriginalURL(URL))

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

type apiError struct {
	Message string `json:"message"`
}

type errorResponse struct {
	Error apiError `json:"error"`
}

func CreateShortURL(
	shortener *services.Shortener,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json")

		rawRequestData, err := io.ReadAll(req.Body)
		if err != nil {
			writeError(res)
			return
		}
		defer req.Body.Close()

		var requestData createShortURLRequest
		err = json.Unmarshal(rawRequestData, &requestData)
		if err != nil {
			writeError(res)
			return
		}

		if len(requestData.URL) == 0 {
			writeError(res)
			return
		}

		shortURL := shortener.Shorten(services.OriginalURL(requestData.URL))

		responseData := createShortURLResponse{
			Result: shortURL,
		}
		rawResponseData, err := json.Marshal(responseData)
		if err != nil {
			writeError(res)
			return
		}

		fmt.Println("test")

		res.WriteHeader(http.StatusCreated)
		res.Write(rawResponseData)
	}
}

func writeError(res http.ResponseWriter) {
	res.WriteHeader(http.StatusBadRequest)

	responseData := errorResponse{
		Error: apiError{
			Message: "bad request",
		},
	}

	rawResponseData, err := json.Marshal(responseData)
	if err == nil {
		res.Write(rawResponseData)
	}
}
