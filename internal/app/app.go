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
	"github.com/aleksandrpnshkn/go-shortener/internal/services/audit"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

const codesLength = 8
const reservedCodesCount = 30

func Run(
	ctx context.Context,
	config *config.Config,
	logger *zap.Logger,
	urlsStorage urls.Storage,
	usersStorage users.Storage,
) error {
	router := chi.NewRouter()

	codeGenerator := services.NewRandomCodeGenerator(codesLength)

	codesReserver := services.NewCodesReserver(
		ctx,
		logger,
		codeGenerator,
		urlsStorage,
		reservedCodesCount,
	)

	auditObservers := []audit.Observer{}

	fileLogger, err := audit.NewFileLogger(config.Audit.File)
	if err == nil {
		defer fileLogger.Close()

		ch := make(chan audit.Event, 300)
		defer close(ch)

		fileObserver := audit.NewFileObserver(ctx, logger, fileLogger, ch)
		auditObservers = append(auditObservers, fileObserver)
	} else {
		logger.Error("failed to create audit file observer", zap.Error(err))
	}

	if config.Audit.URL != "" {
		remoteLogger := audit.NewRemoteLogger(&http.Client{}, config.Audit.URL)
		remoteObserver := audit.NewRemoteObserver(logger, remoteLogger)
		auditObservers = append(auditObservers, remoteObserver)
	}

	shortenedPublisher := audit.NewPublisher(auditObservers)

	shortener := services.NewShortener(
		ctx,
		codesReserver,
		urlsStorage,
		config.PublicBaseURL,
		shortenedPublisher,
	)

	followedPublisher := audit.NewPublisher(auditObservers)

	unshortener := services.NewUnshortener(urlsStorage, followedPublisher)

	auther := services.NewAuther(usersStorage, config.AuthSecretKey)
	deletionBatcher := services.NewDeletionBatcher(ctx, logger, urlsStorage)
	defer deletionBatcher.Close()

	var deleteUserUrlsWg sync.WaitGroup

	router.Use(middlewares.NewLogMiddleware(logger))
	router.Use(compress.NewDecompressMiddleware(logger))
	router.Use(compress.NewCompressMiddleware(logger))

	router.Get("/ping", handlers.PingHandler(ctx, urlsStorage, logger))

	// only login
	router.Group(func(router chi.Router) {
		router.Get("/{code}", handlers.GetURLByCode(auther, unshortener))
	})

	// login or register
	router.Group(func(router chi.Router) {
		router.Use(middlewares.NewAuthMiddleware(logger, auther, true))

		router.Post("/", handlers.CreateShortURLPlain(shortener, logger, auther))
		router.Get("/", handlers.FallbackHandler())

		router.Post("/api/shorten", handlers.CreateShortURL(shortener, logger, auther))
		router.Post("/api/shorten/batch", handlers.CreateShortURLBatch(shortener, logger, auther))
		router.Get("/api/user/urls", handlers.GetUserURLs(shortener, logger, auther))
		router.Delete("/api/user/urls", handlers.DeleteUserURLs(logger, auther, &deleteUserUrlsWg, deletionBatcher))
	})

	if config.EnablePprof {
		logger.Info("enabling pprof routes...")
		router.Mount("/debug", middleware.Profiler())
	}

	logger.Info("running app...")

	err = http.ListenAndServe(config.ServerAddr, router)

	deleteUserUrlsWg.Wait()

	return err
}
