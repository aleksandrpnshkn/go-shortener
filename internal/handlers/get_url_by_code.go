package handlers

import (
	"net/http"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

func GetURLByCode(urlsStorage urls.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain")

		url, err := urlsStorage.Get(req.Context(), types.Code(req.PathValue("code")))
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if url.IsDeleted {
			res.WriteHeader(http.StatusGone)
			return
		}

		res.Header().Add("Location", string(url.OriginalURL))
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}
