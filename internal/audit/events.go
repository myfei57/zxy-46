package audit

import "time"

type Entry struct {
	ID     string
	Source string
	ZoneID string
	Action string
	Detail string
	At     time.Time
}
