package batcher

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
)

type BatchParam any

type BatchExecutor interface {
	GetName() string
	Execute(ctx context.Context, params []BatchParam) error
}

type Batcher struct {
	appCtx        context.Context
	logger        *zap.Logger
	batchExecutor BatchExecutor
	batchSize     int
	inputCh       chan BatchParam
	t             *time.Ticker
}

func (b *Batcher) Add(ctx context.Context, param BatchParam) {
	if len(b.inputCh) >= b.batchSize {
		go func() {
			go b.RunBatch()
		}()
	}

	select {
	case <-ctx.Done():
		return
	case b.inputCh <- param:
		return
	}
}

func (b *Batcher) RunBatch() {
	params := []BatchParam{}

loop:
	for {
		select {
		case param := <-b.inputCh:
			params = append(params, param)
		default:
			break loop
		}
	}

	if len(params) > 0 {
		ctx, cancel := context.WithTimeout(b.appCtx, 5*time.Second)
		defer cancel()

		b.logger.Info("running batch execution...",
			zap.String("name", b.batchExecutor.GetName()),
			zap.Int("batch_size", len(params)),
		)

		err := b.batchExecutor.Execute(ctx, params)
		if err != nil {
			b.logger.Error("error occurred during batch execution",
				zap.String("name", b.batchExecutor.GetName()),
				zap.Error(err),
			)
		}
	}
}

// завершать после остановки хэндлеров
func (b *Batcher) Close() error {
	select {
	case <-b.appCtx.Done():
		close(b.inputCh)
		b.t.Stop()
		return nil
	default:
		return errors.New("app context is not done yet, cannot close channel")
	}
}

func NewBatcher(
	appCtx context.Context,
	logger *zap.Logger,
	batchExecutor BatchExecutor,
	batchSize int,
	batchDelay time.Duration,
) *Batcher {
	ticker := time.NewTicker(batchDelay)

	q := &Batcher{
		appCtx:        appCtx,
		logger:        logger,
		batchExecutor: batchExecutor,
		inputCh:       make(chan BatchParam, batchSize),
		batchSize:     batchSize,
		t:             ticker,
	}

	// регулярно запускать пачки в работу не дожидаясь лимита по размеру
	go func() {
		for {
			select {
			case <-appCtx.Done():
				return
			case <-ticker.C:
				q.RunBatch()
			}
		}
	}()

	return q
}
