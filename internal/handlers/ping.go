package handlers

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
)

// PingHandler - хендлер для проверки состояния веб-сервера.
func PingHandler(
	ctx context.Context,
	storage urls.Storage,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain")

		err := storage.Ping(ctx)
		if err != nil {
			logger.Error("failed to ping sql storage", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
		}

		res.WriteHeader(http.StatusOK)
	}
}
