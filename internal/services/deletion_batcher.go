package services

import (
	"context"
	"time"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
	"go.uber.org/zap"
)

const batchSize = 100

type DeletionBatcher struct {
	logger      *zap.Logger
	ch          chan urls.DeleteCode
	t           *time.Ticker
	urlsStorage urls.Storage
}

type DeleteCode struct {
	Code types.Code
	User users.User
}

func (q *DeletionBatcher) AddCode(ctx context.Context, code types.Code, user users.User) {
	if len(q.ch) >= batchSize {
		q.RunBatch(ctx)
	}

	q.ch <- urls.DeleteCode{
		Code: code,
		User: user,
	}
}

func (q *DeletionBatcher) RunBatch(ctx context.Context) {
	commands := []urls.DeleteCode{}

loop:
	for {
		select {
		case <-ctx.Done():
			return
		case cmd := <-q.ch:
			commands = append(commands, cmd)
		default:
			break loop
		}
	}

	if len(commands) > 0 {
		q.logger.Info("running deletion batch...", zap.Int("commands_size", len(commands)))
		err := q.urlsStorage.DeleteManyByUserID(ctx, commands)
		if (err != nil) {
			q.logger.Error("error occurred during deletion", zap.Error(err))
		}
	}
}

func (q *DeletionBatcher) Close() {
	close(q.ch)
	q.t.Stop()
}

func NewDeletionBatcher(logger *zap.Logger, urlsStorage urls.Storage) *DeletionBatcher {
	ticker := time.NewTicker(time.Second)

	q := &DeletionBatcher{
		ch:          make(chan urls.DeleteCode, batchSize),
		t:           ticker,
		logger:      logger,
		urlsStorage: urlsStorage,
	}

	// регулярно запускать не дожидаясь лимита по размеру
	go func() {
		for {
			<-ticker.C
			q.RunBatch(context.Background())
		}
	}()

	return q
}
