package handlers

import (
	"context"
	"net/http"

	"github.com/aleksandrpnshkn/go-shortener/internal/services"
)

func GetURLByCode(URLsStorage *services.URLsStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain")

		code := services.Code(req.PathValue("code"))
		url, ok := URLsStorage.Get(context.Background(), code)

		if !ok {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.Header().Add("Location", string(url))
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}
