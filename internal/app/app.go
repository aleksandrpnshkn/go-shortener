package app

import (
	"net/http"

	"github.com/aleksandrpnshkn/go-shortener/internal/config"
	"github.com/aleksandrpnshkn/go-shortener/internal/handlers"
	"github.com/aleksandrpnshkn/go-shortener/internal/middlewares"
	"github.com/aleksandrpnshkn/go-shortener/internal/middlewares/compress"
	"github.com/aleksandrpnshkn/go-shortener/internal/services"
	"github.com/aleksandrpnshkn/go-shortener/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const codesLength = 8

func Run(config *config.Config, logger *zap.Logger, store store.Storage) error {
	router := chi.NewRouter()

	codeGenerator := services.NewRandomCodeGenerator(codesLength)
	URLsStorage := services.NewURLsStorage(store)
	shortener := services.NewShortener(
		codeGenerator,
		URLsStorage,
		config.PublicBaseURL,
	)

	router.Use(middlewares.NewLogMiddleware(logger))
	router.Use(compress.DecompressMiddleware)
	router.Use(compress.CompressMiddleware)

	router.Get("/{code}", handlers.GetURLByCode(URLsStorage))
	router.Post("/", handlers.CreateShortURLPlain(shortener))
	router.Get("/", handlers.FallbackHandler())

	router.Post("/api/shorten", handlers.CreateShortURL(shortener))

	logger.Info("Running app...")

	return http.ListenAndServe(config.ServerAddr, router)
}
