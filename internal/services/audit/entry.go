package audit

import "strconv"

type entry struct {
	TimeTS      int64  `json:"ts"`
	Action      string `json:"action"`
	UserID      string `json:"user_id,omitzero"`
	OriginalURL string `json:"url"`
}

func newEntryFromAuditEvent(event Event) entry {
	return entry{
		TimeTS:      event.getTime().Unix(),
		Action:      event.getName(),
		UserID:      strconv.FormatInt(int64(event.getUserID()), 10),
		OriginalURL: string(event.getOriginalURL()),
	}
}
