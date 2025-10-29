package audit

import (
	"context"
	"time"

	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

type Observer interface {
	HandleEvent(ctx context.Context, event Event)
}

type Event interface {
	GetTime() time.Time
	GetName() string
	GetUserID() int64
	GetOriginalURL() types.OriginalURL
}

type Publisher struct {
	observers []Observer
}

func (p *Publisher) Notify(ctx context.Context, event Event) {
	for _, observer := range p.observers {
		observer.HandleEvent(ctx, event)
	}
}

func NewPublisher(observers []Observer) *Publisher {
	return &Publisher{
		observers: observers,
	}
}
