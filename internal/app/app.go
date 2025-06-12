package app

import (
	"math/rand"
	"net/http"
)

type config struct {
	schema   string
	hostname string
}

type application struct {
	config      config
	codesToURLs map[string]string
}

var codeRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

// https://stackoverflow.com/a/31832326
func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = codeRunes[rand.Intn(len(codeRunes))]
	}
	return string(b)
}

func Run() {
	mux := http.NewServeMux()

	config := config{
		schema:   "http",
		hostname: "localhost:8080",
	}

	app := application{
		config:      config,
		codesToURLs: make(map[string]string),
	}

	mux.HandleFunc("GET /{code}", getUrlByCode(app))
	mux.HandleFunc("POST /", createShortUrl(app))
	mux.HandleFunc("/", fallbackHandler())

	err := http.ListenAndServe(app.config.hostname, mux)

	if err != nil {
		panic(err)
	}
}
