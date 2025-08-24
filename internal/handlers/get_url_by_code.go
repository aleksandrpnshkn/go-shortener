package handlers

import (
	"net/http"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
)

func GetURLByCode(urlsStorage urls.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain")

		url, err := urlsStorage.Get(req.Context(), req.PathValue("code"))
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.Header().Add("Location", string(url))
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}
