package demand

func (s *Service) ShedSummary() []string {
	return s.zoneSvc.ShedZones()
}
