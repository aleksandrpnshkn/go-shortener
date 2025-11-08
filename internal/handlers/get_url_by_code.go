package handlers

import (
	"errors"
	"net/http"

	"github.com/aleksandrpnshkn/go-shortener/internal/services"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

// GetURLByCode - хендлер для открытия сокращённого URLа.
func GetURLByCode(auther services.Auther, unshortener *services.Unshortener) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain")

		// nil в юзере это гость. Это ок
		user, _ := auther.FromUserContext(req.Context())

		code := types.Code(req.PathValue("code"))
		originalURL, err := unshortener.Unshorten(req.Context(), code, user)
		if err != nil {
			if errors.Is(err, services.ErrShortURLWasDeleted) {
				res.WriteHeader(http.StatusGone)
				return
			}

			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.Header().Add("Location", string(originalURL))
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}
