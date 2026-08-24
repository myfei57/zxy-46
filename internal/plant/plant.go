package plant

import (
	"errors"
	"sync"

	"bldghvac/internal/audit"
	"bldghvac/internal/chiller"
	"bldghvac/internal/quota"
	"bldghvac/internal/store"
	"bldghvac/internal/weather"
)

type Service struct {
	st       *store.Store
	weather  *weather.Service
	audit    *audit.Service
	quota    *quota.Service
	mu       sync.Mutex
	chillers map[string]*chiller.Unit
	order    []string
	pumps    map[string]chiller.Pump
	setpoint *SetpointState
	arbit    *chiller.LeadRegistry
}

func NewService(st *store.Store, weatherSvc *weather.Service, auditSvc *audit.Service, quotaSvc *quota.Service) *Service {
	s := &Service{
		st:       st,
		weather:  weatherSvc,
		audit:    auditSvc,
		quota:    quotaSvc,
		chillers: map[string]*chiller.Unit{},
		pumps:    map[string]chiller.Pump{},
		setpoint: NewSetpointState(7),
		arbit:    chiller.NewLeadRegistry(),
	}
	var records []chiller.Unit
	if err := st.ReadJSON("plant/chillers.json", &records); err == nil {
		for i := range records {
			unit := records[i]
			s.chillers[unit.ID] = &unit
			s.order = append(s.order, unit.ID)
		}
	}
	var pumpRecords []chiller.Pump
	if err := st.ReadJSON("plant/pumps.json", &pumpRecords); err == nil {
		for _, p := range pumpRecords {
			s.pumps[p.ID] = p
		}
	}
	return s
}

func (s *Service) RegisterChiller(u chiller.Unit) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.chillers[u.ID]; ok {
		return errors.New("chiller already registered")
	}
	unit := u
	s.chillers[unit.ID] = &unit
	s.order = append(s.order, unit.ID)
	return s.persistChillers()
}

func (s *Service) Chiller(id string) (*chiller.Unit, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.chillerLocked(id)
}

func (s *Service) chillerLocked(id string) (*chiller.Unit, error) {
	unit, ok := s.chillers[id]
	if !ok {
		return nil, errors.New("chiller not found")
	}
	return unit, nil
}

func (s *Service) AllChillers() []chiller.Unit {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]chiller.Unit, 0, len(s.order))
	for _, id := range s.order {
		out = append(out, *s.chillers[id])
	}
	return out
}

func (s *Service) Setpoint() *SetpointState {
	return s.setpoint
}

func (s *Service) Baseline() float64 {
	return s.setpoint.Baseline()
}

func (s *Service) persistChillers() error {
	records := make([]chiller.Unit, 0, len(s.order))
	for _, id := range s.order {
		records = append(records, *s.chillers[id])
	}
	return s.st.WriteJSON("plant/chillers.json", records)
}
