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
	local := now.Add(time.Duration(profile.LocalShift) * time.Minute)
	windowStart := effective.Add(-prestartWindowMinutes * time.Minute)
	should := !local.Before(windowStart) && local.Before(effective)
	return PrestartDecision{ShouldStart: should, Effective: effective}
}
