package app

import (
	"io"
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

func createShortURL(app application) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain")

		URL, err := io.ReadAll(req.Body)

		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		defer req.Body.Close()

		var code string
		for codeExists := true; codeExists; {
			code = randStringRunes(8)
			_, codeExists = app.codesToURLs[code]
		}

		app.codesToURLs[code] = string(URL)

		shortURL := app.config.PublicBaseURL + "/" + code

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(shortURL))
	}
}

func fallbackHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain")
		res.WriteHeader(http.StatusBadRequest)
	}
}
