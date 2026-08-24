package vav

func (s *Service) AdjustBaseline(delta float64) float64 {
	s.setpoint.BaselineValue += delta
	return s.setpoint.BaselineValue
}

func (s *Service) Baseline() float64 {
	return s.setpoint.Baseline()
}
