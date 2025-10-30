package audit

import "strconv"

type Entry struct {
	TimeTS      int64  `json:"ts"`
	Action      string `json:"action"`
	UserID      string `json:"user_id,omitzero"`
	OriginalURL string `json:"url"`
}

func NewEntryFromAuditEvent(event Event) Entry {
	return Entry{
		TimeTS:      event.GetTime().Unix(),
		Action:      event.GetName(),
		UserID:      strconv.FormatInt(event.GetUserID(), 10),
		OriginalURL: string(event.GetOriginalURL()),
	}
}
