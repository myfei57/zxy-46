package vav

import (
	"errors"
	"sync"

	"bldghvac/internal/plant"
	"bldghvac/internal/store"
)

type Terminal struct {
	ID         string
	ZoneID     string
	Name       string
	Damper     float64
	ReheatOpen bool
	Flow       float64
}

type Service struct {
	st       *store.Store
	setpoint *plant.SetpointState
	mu       sync.Mutex
	terms    map[string]*Terminal
	order    []string
}

func NewService(st *store.Store, setpoint *plant.SetpointState) *Service {
	s := &Service{
		st: st, setpoint: setpoint,
		terms: map[string]*Terminal{},
	}
	var records []Terminal
	if err := st.ReadJSON("vav/terminals.json", &records); err == nil {
		for i := range records {
			term := records[i]
			s.terms[term.ID] = &term
			s.order = append(s.order, term.ID)
		}
	}
	return s
}

func (s *Service) RegisterTerminal(t Terminal) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.terms[t.ID]; ok {
		return errors.New("terminal already registered")
	}
	term := t
	s.terms[term.ID] = &term
	s.order = append(s.order, term.ID)
	return s.persist()
}

func (s *Service) Terminal(id string) (*Terminal, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	term, ok := s.terms[id]
	if !ok {
		return nil, errors.New("terminal not found")
	}
	return term, nil
}

func (s *Service) AllTerminals() []Terminal {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Terminal, 0, len(s.order))
	for _, id := range s.order {
		out = append(out, *s.terms[id])
	}
	return out
}

func (s *Service) terminalForZone(zoneID string) (*Terminal, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, term := range s.terms {
		if term.ZoneID == zoneID {
			return term, true
		}
	}
	return nil, false
}

func (s *Service) persist() error {
	records := make([]Terminal, 0, len(s.order))
	for _, id := range s.order {
		records = append(records, *s.terms[id])
	}
	return s.st.WriteJSON("vav/terminals.json", records)
}
