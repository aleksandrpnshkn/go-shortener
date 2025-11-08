package audit

import (
	"context"
	"time"

	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

// Observer - подписчик на события аудита.
type Observer interface {
	// HandleEvent обрабатывает событие.
	HandleEvent(ctx context.Context, event Event)
}

// Event - событие аудита.
type Event interface {
	getTime() time.Time
	getName() string
	getUserID() types.UserID
	getOriginalURL() types.OriginalURL
}

// Publisher - генератор событий аудита.
type Publisher struct {
	observers []Observer
}

// Notify уведомляет подписчиков о новом событии аудита.
func (p *Publisher) Notify(ctx context.Context, event Event) {
	for _, observer := range p.observers {
		observer.HandleEvent(ctx, event)
	}
}

// NewPublisher создаёт генератор событий аудита.
func NewPublisher(observers []Observer) *Publisher {
	return &Publisher{
		observers: observers,
	}
}
