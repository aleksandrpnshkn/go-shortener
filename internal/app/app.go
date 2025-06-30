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

type application struct {
	config      *config.Config
	logger      *zap.Logger
	codesToURLs map[string]string
}

func Run(config *config.Config, logger *zap.Logger) {
	router := chi.NewRouter()

	app := application{
		config:      config,
		logger:      logger,
		codesToURLs: make(map[string]string),
	}

	codeGenerator := services.NewCodeGenerator(codesLength)
	fullURLsStorage := services.NewFullURLsStorage()
	shortener := services.NewShortener(
		*codeGenerator,
		*fullURLsStorage,
		app.config.PublicBaseURL,
	)

	router.Use(log.NewRequestMiddleware(logger))

	router.Get("/{code}", getURLByCode(app))
	router.Post("/", handlers.CreateShortURL(shortener))
	router.Get("/", fallbackHandler())

	logger.Info("Running app...")

	err := http.ListenAndServe(app.config.ServerAddr, router)
	if err != nil {
		panic(err)
	}
}
