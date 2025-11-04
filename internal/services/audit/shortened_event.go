package audit

import (
	"time"

	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

type ShortenedEvent struct {
	time        time.Time
	userID      types.UserID
	originalURL types.OriginalURL
}

func (e *ShortenedEvent) GetTime() time.Time {
	return e.time
}

func (e *ShortenedEvent) GetName() string {
	return "shorten"
}

func (e *ShortenedEvent) GetUserID() types.UserID {
	return e.userID
}

func (e *ShortenedEvent) GetOriginalURL() types.OriginalURL {
	return e.originalURL
}

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
