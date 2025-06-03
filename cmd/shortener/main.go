package main

import (
	"io"
	"math/rand"
	"net/http"
)

var codeRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

// https://stackoverflow.com/a/31832326
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = codeRunes[rand.Intn(len(codeRunes))]
	}
	return string(b)
}

func main() {
	schema := "http"
	hostname := "localhost:8080"

	codesToURLs := make(map[string]string)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain")

		if req.Method == http.MethodPost {
			URL, err := io.ReadAll(req.Body)

			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			defer req.Body.Close()

			var code string
			for codeExists := true; codeExists; {
				code = RandStringRunes(8)
				_, codeExists = codesToURLs[code]
			}

			codesToURLs[code] = string(URL)

			shortURL := schema + "://" + hostname + "/" + code

			res.WriteHeader(http.StatusCreated)
			res.Write([]byte(shortURL))
		} else if req.Method == http.MethodGet {
			code := req.URL.Path[1:]
			url, ok := codesToURLs[code]

			if !ok {
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			res.Header().Add("Location", url)
			res.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			res.WriteHeader(http.StatusBadRequest)
		}
	})

	err := http.ListenAndServe(hostname, mux)

	if err != nil {
		panic(err)
	}
}
