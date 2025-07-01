package app

import (
	"net/http"

	"github.com/aleksandrpnshkn/go-shortener/internal/config"
	"github.com/aleksandrpnshkn/go-shortener/internal/handlers"
	"github.com/aleksandrpnshkn/go-shortener/internal/log"
	"github.com/aleksandrpnshkn/go-shortener/internal/services"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const codesLength = 8

func Run(config *config.Config, logger *zap.Logger) {
	router := chi.NewRouter()

	codeGenerator := services.NewRandomCodeGenerator(codesLength)
	fullURLsStorage := services.NewFullURLsStorage()
	shortener := services.NewShortener(
		codeGenerator,
		fullURLsStorage,
		config.PublicBaseURL,
	)

	router.Use(log.NewRequestMiddleware(logger))

	router.Get("/{code}", handlers.GetURLByCode(fullURLsStorage))
	router.Post("/", handlers.CreateShortURLPlain(shortener))
	router.Get("/", handlers.FallbackHandler())

	router.Post("/api/shorten", handlers.CreateShortURL(shortener))

	logger.Info("Running app...")

	err := http.ListenAndServe(config.ServerAddr, router)
	if err != nil {
		panic(err)
	}
}
