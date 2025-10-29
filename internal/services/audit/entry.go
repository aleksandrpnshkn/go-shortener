package audit

import "strconv"

type Entry struct {
	TimeTs      int64  `json:"ts"`
	Action      string `json:"action"`
	UserID      string `json:"user_id"`
	OriginalURL string `json:"url"`
}

func NewEntryFromAuditEvent(event Event) Entry {
	return Entry{
		TimeTs:      event.GetTime().Unix(),
		Action:      event.GetName(),
		UserID:      strconv.FormatInt(event.GetUserID(), 10),
		OriginalURL: string(event.GetOriginalURL()),
	}
}
