package services

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/aleksandrpnshkn/go-shortener/internal/services/batcher"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

// DeleteCode - команда для удаления короткой ссылки из БД
type DeleteCode struct {
	Code   types.Code
	UserID types.UserID
}

// DeleteURLsExecutor удаляет пачку ссылок, которую ему передаёт batcher.Batcher
type DeleteURLsExecutor struct {
	urlsStorage urls.Storage
}

func (e *DeleteURLsExecutor) GetName() string {
	return "delete_orders_executor"
}

func (e *DeleteURLsExecutor) Execute(
	ctx context.Context,
	params []batcher.BatchParam,
) error {
	deleteCommands := make([]urls.DeleteCode, 0, len(params))

	for _, param := range params {
		deleteCode, ok := param.(DeleteCode)
		if !ok {
			return errors.New("passed command is not DeleteCode")
		}
		deleteCommands = append(deleteCommands, urls.DeleteCode{
			Code:   deleteCode.Code,
			UserID: deleteCode.UserID,
		})
	}

	return e.urlsStorage.DeleteManyByUserID(ctx, deleteCommands)
}

func NewDeleteURLsBatcher(
	ctx context.Context,
	logger *zap.Logger,
	urlsStorage urls.Storage,
) *batcher.Batcher {
	deleteURLsExecutor := &DeleteURLsExecutor{
		urlsStorage: urlsStorage,
	}

	batchSize := 100
	batchDelay := 200 * time.Millisecond

	q := batcher.NewBatcher(ctx, logger, deleteURLsExecutor, batchSize, batchDelay)

	return q
}
