package schedule

import (
	"errors"
	"sync"
	"time"

	"bldghvac/internal/store"
)

type Schedule struct {
	ID            string
	ZoneID        string
	Start         time.Time
	OffsetMinutes int
}

const scheduleFile = "schedule/schedules.json"

type Service struct {
	st       *store.Store
	mu       sync.RWMutex
	schedules []Schedule
	offsets  map[string]int
}

func NewService(st *store.Store) *Service {
	s := &Service{st: st, offsets: map[string]int{}}
	var records []Schedule
	if err := st.ReadJSON(scheduleFile, &records); err == nil {
		s.schedules = records
	}
	var offsetRecords map[string]int
	if err := st.ReadJSON("schedule/offsets.json", &offsetRecords); err == nil {
		s.offsets = offsetRecords
	}
	return s
}

func (s *Service) Add(sch Schedule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, existing := range s.schedules {
		if existing.ID == sch.ID {
			return errors.New("schedule already exists")
		}
	}
	s.schedules = append(s.schedules, sch)
	return s.persist()
}

func (s *Service) List() []Schedule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Schedule, len(s.schedules))
	copy(out, s.schedules)
	return out
}

func (s *Service) ForZone(zoneID string) (Schedule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, sch := range s.schedules {
		if sch.ZoneID == zoneID {
			return sch, nil
		}
	}
	return Schedule{}, errors.New("schedule not found")
}

func (s *Service) persist() error {
	if err := s.st.WriteJSON(scheduleFile, s.schedules); err != nil {
		return err
	}
	return s.st.WriteJSON("schedule/offsets.json", s.offsets)
}
