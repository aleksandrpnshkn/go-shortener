package app

import (
	"math/rand"
	"net/http"

	"github.com/aleksandrpnshkn/go-shortener/internal/config"
	"github.com/go-chi/chi/v5"
)

type application struct {
	config      config.Config
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

func Run(config config.Config) {
	router := chi.NewRouter()

	app := application{
		config:      config,
		codesToURLs: make(map[string]string),
	}

	router.Get("/{code}", getURLByCode(app))
	router.Post("/", createShortURL(app))
	router.Get("/", fallbackHandler())

	err := http.ListenAndServe(app.config.ServerAddr, router)
	if err != nil {
		panic(err)
	}
}
