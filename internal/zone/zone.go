package zone

import (
	"errors"
	"sync"

	"bldghvac/internal/store"
)

type Profile struct {
	ZoneID     string
	Setpoint   float64
	Deadband   float64
	MinDamper  float64
	LocalShift int
}

type State struct {
	Setback bool
	Shed    bool
}

const profileFile = "zone/profiles.json"

type Service struct {
	st       *store.Store
	mu       sync.RWMutex
	profiles map[string]Profile
	states   map[string]State
}

func NewService(st *store.Store) *Service {
	s := &Service{st: st, profiles: map[string]Profile{}, states: map[string]State{}}
	var records []Profile
	if err := st.ReadJSON(profileFile, &records); err == nil {
		for _, p := range records {
			s.profiles[p.ZoneID] = p
		}
	}
	return s
}

func (s *Service) UpsertProfile(p Profile) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.profiles[p.ZoneID] = p
	return s.persist()
}

func (s *Service) Profile(zoneID string) (Profile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.profiles[zoneID]
	if !ok {
		return Profile{}, errors.New("zone profile missing")
	}
	return p, nil
}

func (s *Service) AllProfiles() []Profile {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Profile, 0, len(s.profiles))
	for _, p := range s.profiles {
		out = append(out, p)
	}
	return out
}

func (s *Service) persist() error {
	records := make([]Profile, 0, len(s.profiles))
	for _, p := range s.profiles {
		records = append(records, p)
	}
	return s.st.WriteJSON(profileFile, records)
}
