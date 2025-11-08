package audit

import (
	"time"

	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

// ShortenedEvent - событие сокращения URL'а.
type ShortenedEvent struct {
	time        time.Time
	userID      types.UserID
	originalURL types.OriginalURL
}

func (e *ShortenedEvent) getTime() time.Time {
	return e.time
}

func (e *ShortenedEvent) getName() string {
	return "shorten"
}

func (e *ShortenedEvent) getUserID() types.UserID {
	return e.userID
}

func (e *ShortenedEvent) getOriginalURL() types.OriginalURL {
	return e.originalURL
}

// NewShortenedEvent создаёт новое событие сокращения URL'а.
func NewShortenedEvent(
	time time.Time,
	userID types.UserID,
	originalURL types.OriginalURL,
) *ShortenedEvent {
	return &ShortenedEvent{
		time:        time,
		userID:      userID,
		originalURL: originalURL,
	}
}
