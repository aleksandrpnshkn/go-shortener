package handlers

import (
	"net/http"

	"github.com/aleksandrpnshkn/go-shortener/internal/store"
)

func GetURLByCode(urlsStorage store.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain")

		url, ok := urlsStorage.Get(req.Context(), req.PathValue("code"))

		if !ok {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.Header().Add("Location", string(url))
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}
