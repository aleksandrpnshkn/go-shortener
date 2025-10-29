package audit

import (
	"context"

	"go.uber.org/zap"
)

type FileObserver struct {
	logger     *zap.Logger
	fileLogger *FileLogger
}

func (o *FileObserver) HandleEvent(ctx context.Context, event Event) {
	entry := NewEntryFromAuditEvent(event)

	err := o.fileLogger.AddEntry(ctx, entry)
	if err != nil {
		o.logger.Error("observer failed to add entry to file log", zap.Error(err))
	}
}

func NewFileObserver(logger *zap.Logger, fileLogger *FileLogger) *FileObserver {
	return &FileObserver{
		logger:     logger,
		fileLogger: fileLogger,
	}
}
