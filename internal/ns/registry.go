package ns

import (
	"errors"
	"sync"

	"bldghvac/internal/store"
)

const zoneFile = "ns/zones.json"

type Service struct {
	st    *store.Store
	mu    sync.RWMutex
	zones map[string]Zone
	order []string
}

func NewService(st *store.Store) *Service {
	s := &Service{st: st, zones: map[string]Zone{}}
	var records []Zone
	if err := st.ReadJSON(zoneFile, &records); err == nil {
		for _, z := range records {
			s.zones[z.ID] = z
			s.order = append(s.order, z.ID)
		}
	}
	return s
}

func (s *Service) AddZone(z Zone) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.zones[z.ID]; exists {
		return errors.New("zone already exists")
	}
	s.zones[z.ID] = z
	s.order = append(s.order, z.ID)
	return s.persist()
}

func (s *Service) GetZone(id string) (Zone, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	z, ok := s.zones[id]
	if !ok {
		return Zone{}, errors.New("zone not found")
	}
	return z, nil
}

func (s *Service) AllZones() []Zone {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Zone, 0, len(s.order))
	for _, id := range s.order {
		out = append(out, s.zones[id])
	}
	return out
}

func (s *Service) persist() error {
	records := make([]Zone, 0, len(s.order))
	for _, id := range s.order {
		records = append(records, s.zones[id])
	}
	return s.st.WriteJSON(zoneFile, records)
}
