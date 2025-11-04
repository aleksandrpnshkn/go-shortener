package audit

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/aleksandrpnshkn/go-shortener/internal/services/batcher"
)

type EntriesBatcher interface {
	Add(ctx context.Context, entry Entry)
}

type RemoteObserver struct {
	batcher *batcher.Batcher
}

func (o *RemoteObserver) HandleEvent(ctx context.Context, event Event) {
	o.batcher.Add(ctx, NewEntryFromAuditEvent(event))
}

func NewRemoteObserver(
	ctx context.Context,
	logger *zap.Logger,
	remoteLogger *RemoteLogger,
) *RemoteObserver {

	remoteBatchSender := NewRemoteBatchSender(remoteLogger)

	batchSize := 200
	batchDelay := 200 * time.Millisecond

	batcher := batcher.NewBatcher(ctx, logger, remoteBatchSender, batchSize, batchDelay)

	return &RemoteObserver{
		batcher: batcher,
	}
}
