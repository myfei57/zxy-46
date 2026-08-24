package ahu

import (
	"errors"
	"sync"

	"bldghvac/internal/schedule"
	"bldghvac/internal/store"
	"bldghvac/internal/weather"
	"bldghvac/internal/zone"
)

type Unit struct {
	ID          string
	ZoneID      string
	Name        string
	Running     bool
	Prestart    bool
	FreeCooling bool
}

type Service struct {
	st         *store.Store
	zoneSvc    *zone.Service
	schSvc     *schedule.Service
	weatherSvc *weather.Service
	mu         sync.Mutex
	units      map[string]*Unit
	order      []string
}

func NewService(st *store.Store, zoneSvc *zone.Service, schSvc *schedule.Service, weatherSvc *weather.Service) *Service {
	s := &Service{
		st: st, zoneSvc: zoneSvc, schSvc: schSvc, weatherSvc: weatherSvc,
		units: map[string]*Unit{},
	}
	var records []Unit
	if err := st.ReadJSON("ahu/units.json", &records); err == nil {
		for i := range records {
			unit := records[i]
			s.units[unit.ID] = &unit
			s.order = append(s.order, unit.ID)
		}
	}
	return s
}

func (s *Service) Register(u Unit) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.units[u.ID]; ok {
		return errors.New("ahu unit already registered")
	}
	unit := u
	s.units[unit.ID] = &unit
	s.order = append(s.order, unit.ID)
	return s.persist()
}

func (s *Service) Unit(id string) (*Unit, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	unit, ok := s.units[id]
	if !ok {
		return nil, errors.New("ahu unit not found")
	}
	return unit, nil
}

func (s *Service) AllUnits() []Unit {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Unit, 0, len(s.order))
	for _, id := range s.order {
		out = append(out, *s.units[id])
	}
	return out
}

func (s *Service) unitForZone(zoneID string) (*Unit, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, unit := range s.units {
		if unit.ZoneID == zoneID {
			return unit, true
		}
	}
	return nil, false
}

func (s *Service) persist() error {
	records := make([]Unit, 0, len(s.order))
	for _, id := range s.order {
		records = append(records, *s.units[id])
	}
	return s.st.WriteJSON("ahu/units.json", records)
}
