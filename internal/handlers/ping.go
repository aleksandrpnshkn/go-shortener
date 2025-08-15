package handlers

import (
	"context"
	"net/http"

	"github.com/aleksandrpnshkn/go-shortener/internal/store"
)

func PingHandler(ctx context.Context, databaseDSN string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain")

		SQLStorage, err := store.NewSQLStorage(ctx, databaseDSN)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
		}

		err = SQLStorage.Ping(ctx)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
		}

		res.WriteHeader(http.StatusOK)
	}
}
