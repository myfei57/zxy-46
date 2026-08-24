package schedule

import "time"

func IsWorkday(day time.Time) bool {
	weekday := day.Weekday()
	return weekday != time.Saturday && weekday != time.Sunday
}

func NextWorkdayStart(day time.Time, hour, minute int) time.Time {
	candidate := time.Date(day.Year(), day.Month(), day.Day(), hour, minute, 0, 0, day.Location())
	for !IsWorkday(candidate) {
		candidate = candidate.AddDate(0, 0, 1)
	}
	return candidate
}
