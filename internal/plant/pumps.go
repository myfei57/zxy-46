package plant

import (
	"errors"

	"bldghvac/internal/chiller"
)

func (s *Service) RegisterPump(p chiller.Pump) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.pumps == nil {
		s.pumps = map[string]chiller.Pump{}
	}
	s.pumps[p.ID] = p
	return s.persistPumps()
}

func (s *Service) Pump(chillerID string) (chiller.Pump, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.pumpLocked(chillerID)
}

func (s *Service) pumpLocked(chillerID string) (chiller.Pump, error) {
	for _, p := range s.pumps {
		if p.ChillerID == chillerID {
			return p, nil
		}
	}
	return chiller.Pump{}, errors.New("pump not found")
}

func (s *Service) AllPumps() []chiller.Pump {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]chiller.Pump, 0, len(s.pumps))
	for _, p := range s.pumps {
		out = append(out, p)
	}
	return out
}

func (s *Service) persistPumps() error {
	records := make([]chiller.Pump, 0, len(s.pumps))
	for _, p := range s.pumps {
		records = append(records, p)
	}
	return s.st.WriteJSON("plant/pumps.json", records)
}
