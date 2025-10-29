package audit

import (
	"context"
	"strconv"

	"go.uber.org/zap"
)

type FileObserver struct {
	logger     *zap.Logger
	fileLogger *FileLogger
}

func (o *FileObserver) HandleEvent(ctx context.Context, event Event) {
	entry := Entry{
		TimeTs:      event.GetTime().Unix(),
		Action:      event.GetName(),
		UserID:      strconv.FormatInt(event.GetUserID(), 10),
		OriginalURL: string(event.GetOriginalURL()),
	}

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
