package ns

import "github.com/google/uuid"

func NewID(prefix string) string {
	return prefix + "-" + uuid.NewString()[:8]
}

func NewZoneID() string {
	return NewID("zone")
}
