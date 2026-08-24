package schedule

import (
	"time"

	"bldghvac/internal/zone"
)

const prestartWindowMinutes = 15

type PrestartDecision struct {
	ShouldStart bool
	Effective   time.Time
}

func ComputePrestart(profile zone.Profile, offsetMinutes int, baseline time.Time, now time.Time) PrestartDecision {
	effective := baseline.Add(time.Duration(offsetMinutes) * time.Minute)
	windowStart := effective.Add(-prestartWindowMinutes * time.Minute)
	// baseline and now are both expressed in the shared building reference
	// frame. Do not fold the zone LocalShift into the comparison: that would
	// re-interpret the shared baseline as each zone's own local time and
	// stagger prestart by the east/west LocalShift even when the configured
	// offset is zero.
	should := !now.Before(windowStart) && now.Before(effective)
	return PrestartDecision{ShouldStart: should, Effective: effective}
}
