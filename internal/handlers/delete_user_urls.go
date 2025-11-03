package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/aleksandrpnshkn/go-shortener/internal/services"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
	"go.uber.org/zap"
)

func DeleteUserURLs(
	logger *zap.Logger,
	auther services.Auther,
	deleteWg *sync.WaitGroup,
	deletionBatcher *services.Batcher,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json")

		user, err := auther.FromUserContext(req.Context())
		if err != nil {
			logger.Error("failed to get user", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		rawRequestData, err := io.ReadAll(req.Body)
		if err != nil {
			writeBadRequestError(res)
			return
		}
		defer req.Body.Close()

		var codes []types.Code
		err = json.Unmarshal(rawRequestData, &codes)
		if err != nil {
			writeBadRequestError(res)
			return
		}

		deleteWg.Add(1)
		go func() {
			defer deleteWg.Done()

			ctx, cancel := context.WithTimeout(req.Context(), time.Second*10)
			defer cancel()

			for _, code := range codes {
				deleteCodeCommand := services.DeleteCode{
					Code: code,
					User: *user,
				}

				deletionBatcher.Add(ctx, deleteCodeCommand)
			}
		}()

		res.WriteHeader(http.StatusAccepted)
	}
}
