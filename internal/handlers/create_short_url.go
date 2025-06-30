package handlers

import (
	"io"
	"net/http"

	"github.com/aleksandrpnshkn/go-shortener/internal/services"
)

func CreateShortURL(
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

		shortURL := shortener.Shorten(services.FullURL(URL))

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(shortURL))
	}
}
