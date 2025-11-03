package audit

import (
	"context"

	"go.uber.org/zap"
)

type RemoteObserver struct {
	logger       *zap.Logger
	remoteLogger *RemoteLogger
}

func (o *RemoteObserver) HandleEvent(ctx context.Context, event Event) {
	entry := NewEntryFromAuditEvent(event)

	err := o.remoteLogger.SendEntry(ctx, entry)
	if err != nil {
		o.logger.Error("observer failed to send entry to remote log", zap.Error(err))
	}
}

func NewRemoteObserver(logger *zap.Logger, remoteLogger *RemoteLogger) *RemoteObserver {
	return &RemoteObserver{
		logger:       logger,
		remoteLogger: remoteLogger,
	}
}
