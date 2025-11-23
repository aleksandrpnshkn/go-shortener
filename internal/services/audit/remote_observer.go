package audit

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/aleksandrpnshkn/go-shortener/internal/services/batcher"
)

// RemoteObserver - наблюдатель для отправки логов в сторонний сервис.
type RemoteObserver struct {
	batcher *batcher.Batcher
}

// HandleEvent обрабатывает событие.
func (o *RemoteObserver) HandleEvent(ctx context.Context, event Event) {
	o.batcher.Add(ctx, newEntryFromAuditEvent(event))
}

// NewRemoteObserver создаёт нового наблюдателя.
func NewRemoteObserver(
	ctx context.Context,
	logger *zap.Logger,
	auditURL string,
) *RemoteObserver {

	remoteLogger := newRemoteLogger(&http.Client{}, auditURL)

	remoteBatchSender := newRemoteBatchSender(remoteLogger)

	batchSize := 200
	batchDelay := 200 * time.Millisecond

	batcher := batcher.NewBatcher(ctx, logger, remoteBatchSender, batchSize, batchDelay)

	return &RemoteObserver{
		batcher: batcher,
	}
}
