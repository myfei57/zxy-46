package ahu

func (s *Service) VentMode(zoneID string) string {
	unit, ok := s.unitForZone(zoneID)
	if !ok {
		return "unknown"
	}
	if unit.FreeCooling {
		return "free-cooling"
	}
	if unit.Running {
		return "mechanical"
	}
	return "off"
}
