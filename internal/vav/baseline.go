package vav

func (s *Service) AdjustBaseline(delta float64) float64 {
	return s.setpoint.AdjustBaseline(delta)
}

func (s *Service) Baseline() float64 {
	return s.setpoint.Baseline()
}
