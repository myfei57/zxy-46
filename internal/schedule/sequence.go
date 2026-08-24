package schedule

import "time"

func (s *Service) DaySequence(day time.Time) []Schedule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Schedule, 0)
	for _, sch := range s.schedules {
		if IsWorkday(day) == IsWorkday(sch.Start) {
			out = append(out, sch)
		}
	}
	return out
}
