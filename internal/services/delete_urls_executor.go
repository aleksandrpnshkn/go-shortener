package services

import (
	"context"
	"errors"
	"time"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
	"go.uber.org/zap"
)

type DeleteCode struct {
	Code types.Code
	User users.User
}

type DeleteURLsExecutor struct {
	urlsStorage urls.Storage
}

func (e *DeleteURLsExecutor) GetName() string {
	return "delete_orders_executor"
}

func (e *DeleteURLsExecutor) Execute(
	ctx context.Context,
	params []BatchParam,
) error {
	deleteCommands := []urls.DeleteCode{}

	for _, param := range params {
		deleteCode, ok := param.(DeleteCode)
		if !ok {
			return errors.New("passed command is not DeleteCode")
		}
		deleteCommands = append(deleteCommands, urls.DeleteCode{
			Code: deleteCode.Code,
			User: deleteCode.User,
		})
	}

	return e.urlsStorage.DeleteManyByUserID(ctx, deleteCommands)
}

func NewDeletionBatcher(
	ctx context.Context,
	logger *zap.Logger,
	urlsStorage urls.Storage,
) *Batcher {
	deleteURLsExecutor := &DeleteURLsExecutor{
		urlsStorage: urlsStorage,
	}

	batchSize := 200
	batchDelay := 1 * time.Second

	q := NewBatcher(ctx, logger, deleteURLsExecutor, batchSize, batchDelay)

	return q
}
