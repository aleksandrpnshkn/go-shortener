package audit

import (
	"context"

	"go.uber.org/zap"
)

type fileObserver struct {
	logger     *zap.Logger
	ch         chan<- Event
	fileLogger *fileLogger
}

// HandleEvent обрабатывает событие.
func (o *fileObserver) HandleEvent(ctx context.Context, event Event) {
	select {
	case <-ctx.Done():
	case o.ch <- event:
	}
}

// NewFileObserver создаёт нового наблюдателя для записи событий в файл.
func NewFileObserver(
	ctx context.Context,
	logger *zap.Logger,
	fileLogger *fileLogger,
	ch chan Event,
) *fileObserver {
	// писать в файл из одной горутины для избежания data race
	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info("stop file audit logger goroutine")
				return
			case event := <-ch:
				entry := newEntryFromAuditEvent(event)

				err := fileLogger.addEntry(ctx, entry)
				if err != nil {
					logger.Error("observer failed to add entry to file log", zap.Error(err))
				}
			}

		}
	}()

	return &fileObserver{
		logger:     logger,
		ch:         ch,
		fileLogger: fileLogger,
	}
}
