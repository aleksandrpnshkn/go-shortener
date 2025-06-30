package app

import (
	"net/http"
)

func getURLByCode(app application) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain")

		code := req.PathValue("code")
		url, ok := app.codesToURLs[code]

		if !ok {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.Header().Add("Location", url)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func fallbackHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain")
		res.WriteHeader(http.StatusBadRequest)
	}
}
