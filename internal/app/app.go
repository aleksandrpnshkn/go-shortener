package app

import (
	"math/rand"
	"net/http"

	"github.com/aleksandrpnshkn/go-shortener/internal/config"
	"github.com/aleksandrpnshkn/go-shortener/internal/log"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type application struct {
	config      *config.Config
	logger      *zap.Logger
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

func Run(config *config.Config, logger *zap.Logger) {
	router := chi.NewRouter()

	app := application{
		config:      config,
		logger:      logger,
		codesToURLs: make(map[string]string),
	}

	router.Use(log.NewRequestMiddleware(logger))

	router.Get("/{code}", getURLByCode(app))
	router.Post("/", createShortURL(app))
	router.Get("/", fallbackHandler())

	logger.Info("Running app...")

	err := http.ListenAndServe(app.config.ServerAddr, router)
	if err != nil {
		panic(err)
	}
}
