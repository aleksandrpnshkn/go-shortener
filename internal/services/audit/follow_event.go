package audit

import (
	"time"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

type FollowEvent struct {
	time        time.Time
	user        *users.User
	originalURL types.OriginalURL
}

func (e *FollowEvent) GetTime() time.Time {
	return e.time
}

func (e *FollowEvent) GetName() string {
	return "follow"
}

func (e *FollowEvent) GetUserID() int64 {
	if e.user == nil {
		return 0
	}

	return e.user.ID
}

func (e *FollowEvent) GetOriginalURL() types.OriginalURL {
	return e.originalURL
}

func NewFollowEvent(
	time time.Time,
	user *users.User,
	originalURL types.OriginalURL,
) *FollowEvent {
	event := &FollowEvent{
		time:        time,
		user:        user,
		originalURL: originalURL,
	}

	return event
}
