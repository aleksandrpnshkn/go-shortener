package audit

import (
	"time"

	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

type FollowEvent struct {
	time        time.Time
	userID      types.UserID
	originalURL types.OriginalURL
}

func (e *FollowEvent) GetTime() time.Time {
	return e.time
}

func (e *FollowEvent) GetName() string {
	return "follow"
}

func (e *FollowEvent) GetUserID() types.UserID {
	return e.userID
}

func (e *FollowEvent) GetOriginalURL() types.OriginalURL {
	return e.originalURL
}

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
