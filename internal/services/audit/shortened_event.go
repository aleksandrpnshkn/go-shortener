package audit

import (
	"time"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

type ShortenedEvent struct {
	time        time.Time
	user        *users.User
	originalURL types.OriginalURL
}

func (e *ShortenedEvent) GetTime() time.Time {
	return e.time
}

func (e *ShortenedEvent) GetName() string {
	return "shorten"
}

func (e *ShortenedEvent) GetUserID() int64 {
	return e.user.ID
}

func (e *ShortenedEvent) GetOriginalURL() types.OriginalURL {
	return e.originalURL
}

func NewShortenedEvent(
	time time.Time,
	user *users.User,
	originalURL types.OriginalURL,
) *ShortenedEvent {
	return &ShortenedEvent{
		time:        time,
		user:        user,
		originalURL: originalURL,
	}
}
