package audit

import (
	"context"

	"go.uber.org/zap"
)

type FileObserver struct {
	logger     *zap.Logger
	ch         chan<- Event
	fileLogger *FileLogger
}

func (o *FileObserver) HandleEvent(ctx context.Context, event Event) {
	select {
	case <-ctx.Done():
	case o.ch <- event:
	}
}

func NewFileObserver(
	ctx context.Context,
	logger *zap.Logger,
	fileLogger *FileLogger,
	ch chan Event,
) *FileObserver {
	// писать в файл из одной горутины для избежания data race
	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info("stop file audit logger goroutine")
				return
			case event := <-ch:
				entry := NewEntryFromAuditEvent(event)

				err := fileLogger.AddEntry(ctx, entry)
				if err != nil {
					logger.Error("observer failed to add entry to file log", zap.Error(err))
				}
			}

		}
	}()

	return &FileObserver{
		logger:     logger,
		ch:         ch,
		fileLogger: fileLogger,
	}
}
