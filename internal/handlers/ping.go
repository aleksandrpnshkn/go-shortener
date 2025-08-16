package handlers

import (
	"context"
	"net/http"

	"github.com/aleksandrpnshkn/go-shortener/internal/store"
	"go.uber.org/zap"
)

func PingHandler(
	ctx context.Context,
	databaseDSN string,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain")

		sqlStorage, err := store.NewSQLStorage(ctx, databaseDSN)
		if err != nil {
			logger.Error("failed to create sql storage for ping", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
		}

		err = sqlStorage.Ping(ctx)
		if err != nil {
			logger.Error("failed to ping sql storage", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
		}

		res.WriteHeader(http.StatusOK)
	}
}
