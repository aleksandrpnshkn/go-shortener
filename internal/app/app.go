package app

import (
	"context"
	"net/http"
	"sync"

	"github.com/aleksandrpnshkn/go-shortener/internal/config"
	"github.com/aleksandrpnshkn/go-shortener/internal/handlers"
	"github.com/aleksandrpnshkn/go-shortener/internal/middlewares"
	"github.com/aleksandrpnshkn/go-shortener/internal/middlewares/compress"
	"github.com/aleksandrpnshkn/go-shortener/internal/services"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const codesLength = 8

func Run(
	ctx context.Context,
	config *config.Config,
	logger *zap.Logger,
	urlsStorage urls.Storage,
	usersStorage users.Storage,
) error {
	router := chi.NewRouter()

	codeGenerator := services.NewRandomCodeGenerator(codesLength)
	shortener := services.NewShortener(
		codeGenerator,
		urlsStorage,
		config.PublicBaseURL,
	)

	auther := services.NewAuther(usersStorage, config.AuthSecretKey)
	deletionBatcher := services.NewDeletionBatcher(logger, urlsStorage)
	defer deletionBatcher.Close()

	var deleteUserUrlsWg sync.WaitGroup

	router.Use(middlewares.NewLogMiddleware(logger))
	router.Use(compress.NewDecompressMiddleware(logger))
	router.Use(compress.NewCompressMiddleware(logger))
	router.Use(middlewares.NewAuthMiddleware(logger, auther))

	router.Get("/{code}", handlers.GetURLByCode(urlsStorage))
	router.Post("/", handlers.CreateShortURLPlain(shortener, logger, auther))
	router.Get("/", handlers.FallbackHandler())

	router.Post("/api/shorten", handlers.CreateShortURL(shortener, logger, auther))
	router.Post("/api/shorten/batch", handlers.CreateShortURLBatch(shortener, logger, auther))
	router.Get("/api/user/urls", handlers.GetUserURLs(shortener, logger, auther))
	router.Delete("/api/user/urls", handlers.DeleteUserURLs(logger, auther, &deleteUserUrlsWg, deletionBatcher))

	router.Get("/ping", handlers.PingHandler(ctx, urlsStorage, logger))

	logger.Info("Running app...")

	err := http.ListenAndServe(config.ServerAddr, router)

	deleteUserUrlsWg.Wait()

	return err
}
