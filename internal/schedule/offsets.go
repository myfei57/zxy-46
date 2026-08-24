package schedule

func (s *Service) SetOffset(zoneID string, minutes int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.offsets[zoneID] = minutes
	return s.st.WriteJSON("schedule/offsets.json", s.offsets)
}

func (s *Service) Offset(zoneID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.offsets[zoneID]
}
