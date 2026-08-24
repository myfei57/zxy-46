package zone

import (
	"time"
)

type SetbackRecord struct {
	ZoneID        string
	Active        bool
	EnteredAt     time.Time
	ExitedAt      time.Time
	NightBaseline float64
	DayBaseline   float64
}

func (s *Service) EnterSetback(zoneID string, now time.Time, night, day float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	state := s.states[zoneID]
	state.Setback = true
	s.states[zoneID] = state
	return s.st.WriteJSON("zone/setback-"+zoneID+".json", SetbackRecord{
		ZoneID:        zoneID,
		Active:        true,
		EnteredAt:     now,
		NightBaseline: night,
		DayBaseline:   day,
	})
}

func (s *Service) ExitSetback(zoneID string, now time.Time) (SetbackRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var rec SetbackRecord
	if err := s.st.ReadJSON("zone/setback-"+zoneID+".json", &rec); err != nil {
		return SetbackRecord{}, err
	}
	rec.ExitedAt = now
	state := s.states[zoneID]
	state.Setback = false
	s.states[zoneID] = state
	if err := s.st.WriteJSON("zone/setback-"+zoneID+".json", rec); err != nil {
		return SetbackRecord{}, err
	}
	return rec, nil
}

func (s *Service) InSetback(zoneID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.states[zoneID].Setback
}
