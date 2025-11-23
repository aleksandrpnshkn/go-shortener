package audit

import (
	"time"

	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

// FollowEvent - событие открытия ссылки.
type FollowEvent struct {
	time        time.Time
	userID      types.UserID
	originalURL types.OriginalURL
}

func (e *FollowEvent) getTime() time.Time {
	return e.time
}

func (e *FollowEvent) getName() string {
	return "follow"
}

func (e *FollowEvent) getUserID() types.UserID {
	return e.userID
}

func (e *FollowEvent) getOriginalURL() types.OriginalURL {
	return e.originalURL
}

// NewShortenedEvent создаёт новое событие открытия ссылки.
func NewFollowEvent(
	time time.Time,
	userID types.UserID,
	originalURL types.OriginalURL,
) *FollowEvent {
	event := &FollowEvent{
		time:        time,
		userID:      userID,
		originalURL: originalURL,
	}

	return event
}
