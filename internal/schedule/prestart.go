package schedule

import "time"

const prestartWindowMinutes = 15

type PrestartDecision struct {
	ShouldStart bool
	Effective   time.Time
}

func ComputePrestart(offsetMinutes int, baseline time.Time, now time.Time) PrestartDecision {
	effective := baseline.Add(time.Duration(offsetMinutes) * time.Minute)
	windowStart := effective.Add(-prestartWindowMinutes * time.Minute)
	should := !now.Before(windowStart) && now.Before(effective)
	return PrestartDecision{ShouldStart: should, Effective: effective}
}
