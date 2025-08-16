package app

import (
	"context"
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

func Run(ctx context.Context, config *config.Config, logger *zap.Logger, storage store.Storage) error {
	router := chi.NewRouter()

	codeGenerator := services.NewRandomCodeGenerator(codesLength)
	shortener := services.NewShortener(
		codeGenerator,
		storage,
		config.PublicBaseURL,
	)

	router.Use(middlewares.NewLogMiddleware(logger))
	router.Use(compress.NewDecompressMiddleware(logger))
	router.Use(compress.NewCompressMiddleware(logger))

	router.Get("/{code}", handlers.GetURLByCode(storage))
	router.Post("/", handlers.CreateShortURLPlain(shortener, logger))
	router.Get("/", handlers.FallbackHandler())

	router.Post("/api/shorten", handlers.CreateShortURL(shortener, logger))
	router.Post("/api/shorten/batch", handlers.CreateShortURLBatch(shortener, logger))

	router.Get("/ping", handlers.PingHandler(ctx, storage, logger))

	logger.Info("Running app...")

	return http.ListenAndServe(config.ServerAddr, router)
}
