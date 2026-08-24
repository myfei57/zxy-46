package zone

func (s *Service) ApplyShed(zoneID string, shed bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	state := s.states[zoneID]
	state.Shed = shed
	s.states[zoneID] = state
	return s.st.WriteJSON("zone/shed-"+zoneID+".json", state)
}

func (s *Service) ShedState(zoneID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.states[zoneID].Shed
}

func (s *Service) ShedZones() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, 0)
	for id, state := range s.states {
		if state.Shed {
			out = append(out, id)
		}
	}
	return out
}
