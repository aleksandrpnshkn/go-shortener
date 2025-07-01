package app

import (
	"net/http"

	"github.com/aleksandrpnshkn/go-shortener/internal/services"
)

func getURLByCode(fullURLsStorage *services.FullURLsStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain")

		code := services.Code(req.PathValue("code"))
		url, ok := fullURLsStorage.Get(code)

		if !ok {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.Header().Add("Location", string(url))
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func fallbackHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain")
		res.WriteHeader(http.StatusBadRequest)
	}
}
